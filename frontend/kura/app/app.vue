<script setup lang="ts">
// トークンを取得
const token = useCookie("token")
// const token = "1";
</script>

<!--------------------->

<template>
    <div class="app-container">

        <Header></Header>
        
        <main>
            <NuxtPage />
        </main>
    </div>
</template>

<!--------------------->

<style scoped>
div.app-container {
    display: flex;
    flex-direction: column;
    min-height: 100dvh;
}

main {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 2rem;
    flex: 1;

    padding: 1rem;
}
</style>