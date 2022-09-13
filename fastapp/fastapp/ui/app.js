window.addEventListener("load", () => {
  const input = document.getElementById("inputFile");
  const filename = document.getElementById("filename");
  const upload = document.getElementById("upload");

  input.addEventListener("change", () => {
    const inputFile = document.querySelector("input[type=file]").files[0];

    filename.innerText = inputFile.name;
    upload.removeAttribute("disabled");

    clearProgress();
  });

  upload.addEventListener("click", postFile);
});

function postFile(e) {
  e.preventDefault();
  const inputFile = document.querySelector("input[type=file]").files[0];
  const progressBar = document.getElementById("progress-bar");
  const pending = document.getElementById("pending");
  pending.hidden = false;
  var formdata = new FormData();

  formdata.append("file", inputFile);

  var request = new XMLHttpRequest();

  request.upload.addEventListener("progress", function (e) {
    var fileSize = inputFile.size;

    if (e.loaded <= fileSize) {
      var percent = Math.round((e.loaded / fileSize) * 100);
      progressBar.style.width = percent + "%";
      progressBar.innerHTML = percent + "%";
    }

    if (e.loaded == e.total) {
      progressBar.style.width = "100%";
      progressBar.innerHTML = "100%";
    }
  });

  request.onreadystatechange = () => {
    if (request.readyState === 4) {
      const success = document.getElementById("success");
      const fail = document.getElementById("fail");
      success.hidden = true;
      fail.hidden = true;
      pending.hidden = true;

      try {
        const response = JSON.parse(request.response);
        console.log(response);
        if (response.success) {
          success.hidden = false;
        } else {
          fail.hidden = false;
        }
      } catch (error) {
        fail.hidden = false;
      }
    }
  };

  request.open("post", "/upload-hyperpack");
  request.timeout = 45000;
  request.send(formdata);
}

function clearProgress() {
  const success = document.getElementById("success");
  const pending = document.getElementById("success");
  const fail = document.getElementById("fail");
  const progressBar = document.getElementById("progress-bar");

  success.hidden = true;
  pending.hidden = true;
  fail.hidden = true;
  progressBar.style.width = "0px";
  progressBar.innerHTML = "";
}
