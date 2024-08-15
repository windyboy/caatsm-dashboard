/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./**/*.html", 
    "./**/*.templ", 
    "./**/*.go",
    "./node_modules/flowbite/**/*.js"
  ],
  plugins: [
    require('flowbite/plugin')
],
  safelist: [],
};
