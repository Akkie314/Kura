<script setup lang="ts">
const config = useRuntimeConfig()
const token = useCookie("token")

const { state, uploadFile, reset } = useChunkUpload();

function onChange(file: File) {
    if (token.value) {
        uploadFile(file, token.value)
    }
}
</script>

<!--------------------->

<template>
    <!-- トークンがある人 -->
    <template v-if="token != null">
        <template v-if="state.status === 'idle'">
            <UploadBox @change="onChange"></UploadBox>
        </template>

        <template v-else-if="state.status === 'uploading'">
            <p>{{ state.message }}</p>
            <ProgressBar :value="state.progress"></ProgressBar>
        </template>

        <template v-else-if="state.status === 'complete'">
            <p>{{ state.message }}</p>
            <AppButton :href="`${config.public.baseUrl}/file/${state.fileId}/`">アップロードしたファイルのページへ移動</AppButton>
        </template>
    </template>

    <template v-else-if="state.status === 'error'">
        <div class="status-box error">
            <p class="error-title">⚠️ アップロードに失敗しました</p>
            <p class="error-message">{{ state.message }}</p>

            <div class="action-buttons">
                <AppButton @click="reset">最初からやり直す</AppButton>
            </div>
        </div>
    </template>

    <!-- トークンがない人 -->
    <template v-else>
        <LoginButton></LoginButton>
    </template>
</template>

<!--------------------->

<style scoped></style>