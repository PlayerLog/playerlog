/** @type {import('tailwindcss').Config} */
module.exports = {
  // content: [
  //            "./cmd/web/**/*.html", "./cmd/web/**/*.templ",
  // ],
  content: [
    "./cmd/web/**/*.html",
    "./cmd/web/**/*.templ",
    "C:/Users/Gurrag/go/pkg/mod/github.com/axzilla/templui@*/**/*.{go,templ}",
  ],
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px",
      },
    },
    extend: {},
  },
  plugins: [],
};
