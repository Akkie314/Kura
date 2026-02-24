<script setup lang="ts">
import App from '~/app.vue';
const config = useRuntimeConfig()
</script>

<!--------------------->

<template>
    <h1>Kura</h1>
    <p>KuraはDiscord上で10MB以上の動画やファイルを共有するためのWebサービスです. ファイルをアップロードした後, そのファイルページのURLをDiscordで送信することで, 動画や画像であればDiscord上でも表示されます.</p>

    <AppButton :href="`${config.public.baseUrl}/upload`">アップロードページに移動</AppButton>
</template>

<!--------------------->

<style scoped>
p {
    max-width: 600px;
}
</style>