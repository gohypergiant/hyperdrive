const input = document.getElementById("inputFile");
const filename = document.getElementById("filename");
const upload = document.getElementById("upload");

input.addEventListener("change", () => {
  const inputFile = document.querySelector("input[type=file]").files[0];

  filename.innerText = inputFile.name;
  upload.removeAttribute("disabled");
});
