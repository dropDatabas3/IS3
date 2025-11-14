import { defineConfig } from "cypress";

export default defineConfig({
  e2e: {
    baseUrl: "https://front-tp05-qa-palacio-nallar.azurewebsites.net", // ajustar si preferís otra URL
    setupNodeEvents(on, config) {
      // implement node event listeners here (screenshots, videos, etc.)
      return config;
    },
    reporter: "junit",
    reporterOptions: {
      mochaFile: "cypress/results/results-[hash].xml",
      toConsole: true,
    },
  },
  experimentalStudio: true,
});
