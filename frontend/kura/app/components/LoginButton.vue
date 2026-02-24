<script setup lang="ts">
const config = useRuntimeConfig()

// DiscordログインのURL
const backendLoginUrl = `${config.public.baseUrl}/api/auth/login`;
</script>

<template>
    <a :href="backendLoginUrl" class="btn">
        Discordでログイン
    </a>
</template>

<style scoped>
a.btn {
    display: flex;
    padding: 0.75rem 1.5rem;
    justify-content: center;
    align-items: center;
    gap: 10px;

    border-radius: 0.5rem;
    background: var(--color-bg-blue);

    box-shadow: var(--box-shadow-value);

    cursor: pointer;

    transition: background 0.15s ease-out;
}

a:link,
a:visited {
    color: var(--color-text-white);
    text-decoration: none;
}

a.btn:hover,
a.btn:active,
a.btn:focus,
a.btn:has(input:focus) {
    background: var(--color-bg-dark-blue);
}

</style>