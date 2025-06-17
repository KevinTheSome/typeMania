<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from "vue";
import { generate } from "random-words";
import axios from "axios";

import DifficultySelect from "./components/DifficultySelect.vue";
import Navbar from "./components/Navbar.vue";
import { PlayerStats } from "./types.ts";

const gameActive = ref(false);
const gameFinished = ref(false);
const apiURL = import.meta.env.VITE_API_URL;

const words = ref<string[]>([]);
const wordCount = ref(0);
const currentWordIndex = ref(0);
const userInput = ref("");

const timer = ref(0);
const timerInterval = ref<number | undefined>(undefined);
const errors = ref(0);
const totalCharsTyped = ref(0);

const pStat = ref({
  name: null,
  wpm: 0,
  errors: 0,
  characters: 0,
  accuracy: 0,
  wordCount: 0,
  time: 0,
});
const wpm = computed(() => {
  if (timer.value === 0) return 0;
  const wordsTyped = (totalCharsTyped.value - errors.value) / 5;
  const minutes = timer.value / 60;
  return Math.round(wordsTyped / minutes);
});

const accuracy = computed(() => {
  if (totalCharsTyped.value === 0) return 100;
  return Math.round(
    ((totalCharsTyped.value - errors.value) / totalCharsTyped.value) * 100
  );
});

function generateWords() {
  words.value = generate({
    exactly: wordCount.value,
    minLength: 3,
  }) as string[];
}

function resetGame() {
  gameActive.value = false;
  gameFinished.value = false;
  clearInterval(timerInterval.value);
  timer.value = 0;
  currentWordIndex.value = 0;
  userInput.value = "";
  errors.value = 0;
  totalCharsTyped.value = 0;
  generateWords();
}

function publishScore() {
  pStat.value = {
    name: pStat.value.name,
    wpm: wpm.value,
    errors: errors.value,
    characters: totalCharsTyped.value,
    accuracy: accuracy.value,
    time: timer.value,
    wordCount: wordCount.value,
  };
  if (!pStat.value.name) return;

  axios.post(apiURL, pStat.value);

  resetGame();
}

function handleKeyDown(evt: KeyboardEvent) {
  if (gameFinished.value) return;

  if (!gameActive.value) {
    gameActive.value = true;
    startTimer();
  }

  const currentWord = words.value[currentWordIndex.value];

  if (evt.key === " ") {
    evt.preventDefault();
    if (userInput.value.trim() === currentWord) {
      totalCharsTyped.value++;
      currentWordIndex.value++;
      userInput.value = "";

      if (currentWordIndex.value >= words.value.length) {
        finishGame();
      }
    }
  } else if (evt.key === "Backspace") {
    userInput.value = userInput.value.slice(0, -1);
  } else if (evt.key.length === 1 && !evt.ctrlKey && !evt.metaKey) {
    totalCharsTyped.value++;
    userInput.value += evt.key;
    if (
      userInput.value.length > 0 &&
      !currentWord.startsWith(userInput.value)
    ) {
      errors.value++;
    }
  }
}

function startTimer() {
  timerInterval.value = window.setInterval(() => {
    timer.value++;
  }, 1000);
}

function finishGame() {
  gameActive.value = false;
  gameFinished.value = true;
  clearInterval(timerInterval.value);
}

onMounted(() => {
  resetGame();
  window.addEventListener("keydown", handleKeyDown);
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleKeyDown);
  clearInterval(timerInterval.value);
});

function getCharClass(wordIdx: number, charIdx: number, char: string): string {
  if (wordIdx < currentWordIndex.value) {
    return "text-green-500";
  }
  if (wordIdx === currentWordIndex.value) {
    if (charIdx < userInput.value.length) {
      return userInput.value[charIdx] === char
        ? "text-green-500"
        : "text-red-500 underline";
    }
  }
  return "text-gray-500";
}

function diffCallback(difficulty: string) {
  wordCount.value = parseInt(difficulty);
  resetGame();
}
</script>

<template>
  <section
    class="flex flex-col items-center bg-gray-800 text-white min-h-screen font-mono"
  >
    <div class="container mx-auto p-8 flex flex-col min-h-screen items-center">
      <h1 class="text-4xl font-bold mb-8 text-yellow-400">TypeMania</h1>
      <DifficultySelect @difficulty-event="diffCallback" />
      <div class="grid grid-cols-3 gap-8 text-2xl mb-8 w-full max-w-2xl">
        <div class="text-center">
          <span class="text-yellow-400 font-bold">{{ wpm }}</span> WPM
        </div>
        <div class="text-center">
          <span class="text-yellow-400 font-bold">{{ accuracy }}%</span>
          Accuracy
        </div>
        <div class="text-center">
          <span class="text-yellow-400 font-bold">{{ timer }}</span
          >s
        </div>
      </div>

      <div class="bg-gray-700 p-6 rounded-lg text-2xl w-full h-full mb-8">
        <span
          v-for="(word, wordIdx) in words"
          :key="wordIdx"
          class="p-2 inline-block"
          :class="{
            'bg-gray-600 rounded px-1': wordIdx === currentWordIndex,
          }"
        >
          <span
            v-for="(char, charIdx) in word"
            :key="charIdx"
            :class="getCharClass(wordIdx, charIdx, char)"
          >
            {{ char }}
          </span>
        </span>
      </div>

      <div class="text-xl text-gray-400 h-8">
        Current input: {{ userInput }}
      </div>

      <div v-if="gameFinished" class="text-center p-8 bg-gray-700 rounded-lg">
        <h2 class="text-3xl text-yellow-400 mb-4">Test Complete!</h2>
        <p class="text-xl mb-6">
          You achieved <span class="font-bold">{{ wpm }} WPM</span> with
          <span class="font-bold">{{ accuracy }}%</span> accuracy.
        </p>
        <button
          @click="resetGame"
          class="bg-yellow-500 text-gray-900 font-bold py-3 px-6 rounded-lg hover:bg-yellow-400 transition-colors m-2"
        >
          Try Again
        </button>
        <div>
          <input
            type="text"
            v-model="pStat.name"
            placeholder="Enter your name"
            class="bg-gray-700 text-white p-2 rounded-lg m-2 border-2"
            required
          />
          <p v-if="!pStat.name" class="text-yellow-400 text-2xl font-bold mb-4">
            Enter a name
          </p>
          <button
            @click="publishScore"
            class="bg-yellow-500 text-gray-900 font-bold py-3 px-6 rounded-lg hover:bg-yellow-400 transition-colors m-2"
          >
            Publish Score
          </button>
        </div>
      </div>
    </div>
  </section>
</template>
