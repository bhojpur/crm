<template>
  <v-main>
    <div style="display: grid;grid-template-columns: 200px auto;grid-column-gap: 40px;">
      <aside class="d-flex flex-column">
        <v-btn :to="{name:'legal'}" exact>legal</v-btn>
        <v-btn :to="{name:'legal.use-cookie'}">use-cookie</v-btn>
        <v-btn :to="{name:'legal.privacy-policy'}">privacy-policy</v-btn>
        <v-btn :to="{name:'legal.terms-of-use'}">terms-of-use</v-btn>
      </aside>
      <router-view />
    </div>
  </v-main>
</template>
