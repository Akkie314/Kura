<script setup lang="ts">
// ※環境に合わせて適宜変更してください
const config = useRuntimeConfig()
const API_BASE = config.public.baseUrl
const SITE_URL = config.public.baseUrl

const route = useRoute()
const fileID = route.params.id

// バックエンドからメタデータを取得する
const { data: meta } = await useFetch(`${API_BASE}/api/files/${fileID}/meta`, {
    baseURL: API_BASE,
    method: 'GET',
})

// データが取れなかった場合 (404)
if (!meta.value) {
    throw createError({
        statusCode: 404,
        statusMessage: 'Not Found',
    })
}

// ファイルタイプ判定とURL生成
const mime = meta.value.mime_type || 'application/octet-stream'
// API経由ではなく、Nginxが直接返すファイルURLを使うのが高速でベスト
const rawURL = `${SITE_URL}/api/files/${fileID}`
const filename = meta.value.original_name || "無名のファイル"

const width = meta.value.width || -1
const height = meta.value.height || -1


// ▼▼▼ OGPの設定出し分け ▼▼▼
if (mime.startsWith('image/')) {
    // 【画像の場合】
    useImageSeoMeta(filename, rawURL, mime, width, height)

} else if (mime.startsWith('video/')) {
    // 【動画の場合】
    let thumbnail_url = 'https://placehold.co/1280x720/000000/FFFFFF/png?text=Play+Video'
    
    // DBでサムネイルがある(true)と判定されている場合
    if (meta.value.has_thumbnail) {
        thumbnail_url = `${SITE_URL}/api/files/${fileID}/thumb`
    }

    useVideoSeoMeta(filename, thumbnail_url, rawURL, mime, width, height)

} else {
    // 【その他のファイル（zipなど）】
    useSeoMeta({
        ogTitle: filename,
        ogSiteName: "Kura Uploader",
        ogType: "website",
        // 汎用アイコンなどを指定
        ogImage: 'https://placehold.co/600x400/000000/FFFFFF/png?text=File',
        twitterCard: "summary",
    })
}

// --- 以下、OGP設定用のヘルパー関数 ---

// 画像用OGP（Twitterでデカく表示する）
function useImageSeoMeta(
    title: string,
    image_url: string,
    mime_type: string,
    width: number,
    height: number,
) {
    useSeoMeta({
        ogTitle: title,
        ogSiteName: "Kura Uploader",

        ogImage: image_url,
        ogImageSecureUrl: image_url,
        ogImageType: mime_type,
        ogImageAlt: title,
        ogImageWidth: width > 0 ? width : undefined,
        ogImageHeight: height > 0 ? height : undefined,



        ogType: "website",

        twitterCard: "summary_large_image",
        twitterImage: image_url,
    })
}

// 動画用OGP（DiscordやTwitterで直接再生できるようにする）
function useVideoSeoMeta(
    title: string, 
    thumbnail_url: string,
    video_url: string,
    mime_type: string,
    width: number,
    height: number,
) {
    const desc = "Kura Uploaderで共有された動画ファイルです。";

    useSeoMeta({
        ogTitle: title,
        ogDescription: desc,
        ogSiteName: "Kura Uploader",
        ogImage: thumbnail_url,

        ogType: 'video.other',
        // この og:video があれば Discord は勝手に動画プレイヤーを展開してくれます
        ogVideo: video_url,
        ogVideoSecureUrl: video_url,
        ogVideoType: mime_type,

        ogVideoWidth: width > 0 ? width : undefined,
        ogVideoHeight: height > 0 ? height : undefined,



        // ▼▼▼ ここを変更 ▼▼▼
        twitterCard: "summary_large_image", // "player" から変更
        twitterDescription: desc,
        twitterImage: thumbnail_url,
        // twitterPlayer: video_url, // 💥 これは仕様違反になるので必ず削除する！
    })
}
</script>

<template>
    <div class="container">
        <h1>{{ meta?.original_name }}</h1>

        <div class="preview-area">
            <img v-if="mime.startsWith('image/')" :src="rawURL" alt="Preview Image" class="preview-media" />
            <video v-else-if="mime.startsWith('video/')" :src="rawURL" controls class="preview-media"></video>
            <div v-else>
                <p>このファイルはプレビューできません。</p>
                <AppButton :href="rawURL" :download="filename">ダウンロードする</AppButton>
            </div>
        </div>

    </div>
    <AppButton :href="`${config.public.baseUrl}/upload/`">他のファイルをアップロードする</AppButton>
</template>

<style scoped>
.container {
    max-width: 800px;
    margin: 40px auto;
    text-align: center;
    font-family: sans-serif;
}
.preview-area {
    margin-top: 20px;
}
.preview-media {
    max-width: 100%;
    max-height: 70vh;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
</style>