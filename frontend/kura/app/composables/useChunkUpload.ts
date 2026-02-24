import { ref } from "vue";

interface UploadState {
    progress: number,
    status: "idle" | "uploading" | "complete" | "error",
    message: string,
    fileId: string | null,
}

export function useChunkUpload() {
    const state = ref<UploadState>({
        progress: 0,
        status: "idle",
        message: "",
        fileId: null,
    });

    const CHUNK_SIZE = 1024 * 1024 * 50; // 50MB
    const MAX_RETRIES = 3; // 最大リトライ回数
    const RETRY_DELAY = 2000; // リトライ間の基本待機時間(ms)

    // 指定時間待機するヘルパー
    const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

    async function uploadFile(file: File, token: string) {
        state.value.status = "uploading";
        state.value.message = "アップロード準備中...";
        state.value.progress = 0;

        const authHeader = { Authorization: `Bearer ${token}` };

        try {
            const totalChunks = Math.ceil(file.size / CHUNK_SIZE);

            // 1. セッション開始
            const startRes = await $fetch<{ upload_id: string }>("/api/upload/session", {
                method: "POST",
                headers: authHeader,
                body: {
                    file_type: file.type || 'application/octet-stream',
                    file_name: file.name,
                    file_size: file.size,
                    total_chunk: totalChunks,
                }
            });

            const uploadId = startRes.upload_id;

            // 2. チャンクごとのループ
            for (let i = 0; i < totalChunks; i++) {
                const start = i * CHUNK_SIZE;
                const end = Math.min(start + CHUNK_SIZE, file.size);
                const chunk = file.slice(start, end);

                let success = false;
                let attempt = 0;

                // --- チャンク単位のリトライループ ---
                while (!success && attempt < MAX_RETRIES) {
                    try {
                        await $fetch(`/api/upload/${uploadId}`, {
                            method: "PATCH",
                            headers: {
                                ...authHeader,
                                "Content-Type": "application/octet-stream",
                                "X-Chunk-Index": i.toString(),
                                "X-Chunk-Offset": start.toString(),
                            },
                            body: chunk,
                        });
                        success = true; // 成功したらループ脱出
                    } catch (err) {
                        attempt++;
                        if (attempt >= MAX_RETRIES) throw err; // リトライ上限到達

                        const delay = RETRY_DELAY * attempt; // 失敗するごとに待ち時間を増やす
                        state.value.message = `接続に失敗しました。再試行中... (${attempt}/${MAX_RETRIES})`;
                        await sleep(delay);
                    }
                }
                // ------------------------------

                const percent = Math.floor((i + 1) / totalChunks * 100);
                state.value.progress = percent;
                state.value.message = `アップロード中... ${percent}%`;
            }

            // 3. 完了通知
            state.value.message = "サーバー側で結合中...";
            const completeRes = await $fetch<{ file_id: string }>(`/api/upload/${uploadId}/complete`, {
                method: "POST",
                headers: authHeader,
            });

            state.value.status = "complete";
            state.value.progress = 100;
            state.value.message = "アップロード完了";
            state.value.fileId = completeRes.file_id;

        } catch (error: any) {
            console.error("Upload process failed:", error);
            state.value.status = "error";
            state.value.message = "アップロードが中断されました。ネットワークを確認してください。";
        }
    }

    function reset() {
        state.value = {
            progress: 0,
            status: "idle",
            message: "",
            fileId: null,
        };
    }

    return {
        state,
        uploadFile,
        reset,
    }
}