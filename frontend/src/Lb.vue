<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import Navbar from "./components/Navbar.vue";
import axios from "axios";
import type { PlayerStats } from "./types.ts"; // Ensure this import is correct

const apiURL = import.meta.env.VITE_API_URL;

const playerStats = ref<PlayerStats[]>([]);
const filterdPlayerStats = ref<PlayerStats[]>([]);
const search = ref("");

onMounted(() => {
  axios
    .get(apiURL)
    .then((res) => {
      playerStats.value = res.data;
    })
    .then((res) => {
      console.log(playerStats.value);
    })
    .catch((err) => {
      console.error(err);
    });
});

watch(search, () => {
  filterdPlayerStats.value = playerStats.value.filter((player) =>
    player.Name.toLowerCase().includes(search.value.toLowerCase())
  );
});
</script>
<template>
  <div class="min-h-screen bg-gray-900 text-gray-100 font-mono p-4 sm:p-8">
    <div class="max-w-7xl mx-auto pt-12 md:pt-24 flex flex-col items-center">
      <h1 class="text-5xl font-extrabold text-yellow-400 mb-8">Leaderboard</h1>

      <input
        type="text"
        placeholder="Search by name..."
        class="mt-4 mb-12 px-6 py-3 bg-gray-800 text-yellow-200 border-2 border-gray-700 rounded-full focus:border-yellow-500 outline-none text-xl w-full max-w-md placeholder-gray-500 transition-colors duration-300 ease-in-out"
        v-model="search"
      />

      <div
        class="grid grid-cols-1 md:grid-cols-1 lg:grid-cols-1 grid-rows-auto gap-4 w-full max-w-6xl"
      >
        <template v-if="filterdPlayerStats.length > 0">
          <div
            v-for="player in filterdPlayerStats"
            :key="player.Id"
            class="bg-gray-800 p-8 rounded-lg shadow-lg border border-gray-700 hover:border-yellow-600 transition-all duration-300 transform hover:scale-105"
          >
            <p class="text-3xl font-bold text-yellow-300 mb-3 break-words">
              {{ player.Name }}
            </p>
            <p class="text-xl text-gray-200 mb-2">
              <span class="font-semibold text-yellow-400">WPM:</span>
              {{ player.Wpm }}
            </p>
            <p class="text-xl text-gray-200">
              <span class="font-semibold text-yellow-400">Accuracy:</span>
              {{ player.Accuracy }}%
            </p>
            <p class="text-xl text-gray-200">
              <span class="font-semibold text-yellow-400">Word Count:</span>
              {{ player.WordCount }}%
            </p>
            <p class="text-xl text-gray-200">
              <span class="font-semibold text-yellow-400">Errors:</span>
              {{ player.Errors }}%
            </p>
            <p class="text-xl text-gray-200">
              <span class="font-semibold text-yellow-400"
                >Time in secends:</span
              >
              {{ player.Time }}
            </p>
          </div>
        </template>

        <template v-else-if="search && filterdPlayerStats.length === 0">
          <p class="col-span-full text-center text-2xl text-gray-400">
            No players found matching "{{ search }}".
          </p>
        </template>

        <template v-else>
          <div
            v-for="player in playerStats"
            :key="player.Id"
            class="bg-gray-800 p-8 rounded-lg shadow-lg border border-gray-700 hover:border-yellow-600 transition-all duration-300 transform hover:scale-105"
          >
            <p class="text-3xl font-bold text-yellow-300 mb-3 break-words">
              {{ player.Name }}
            </p>
            <p class="text-xl text-gray-200 mb-2">
              <span class="font-semibold text-yellow-400">WPM:</span>
              {{ player.Wpm }}
            </p>
            <p class="text-xl text-gray-200">
              <span class="font-semibold text-yellow-400">Accuracy:</span>
              {{ player.Accuracy }}%
            </p>
            <p class="text-xl text-gray-200">
              <span class="font-semibold text-yellow-400">Word Count:</span>
              {{ player.WordCount }}
            </p>
            <p class="text-xl text-gray-200">
              <span class="font-semibold text-yellow-400">Errors:</span>
              {{ player.Errors }}
            </p>
            <p class="text-xl text-gray-200">
              <span class="font-semibold text-yellow-400"
                >Time in secends:</span
              >
              {{ player.Time }}
            </p>
          </div>
        </template>
      </div>

      <p v-if="playerStats.length === 0" class="text-2xl text-gray-400 mt-12">
        The leaderboard is currently empty. Play a game to post your score!
      </p>
    </div>
  </div>
</template>
