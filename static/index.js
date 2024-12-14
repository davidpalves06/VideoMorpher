const fileForm = document.getElementById("fileForm");
const fileUploadInput = document.getElementById("fileUpload");
const sizeUnits = ['Bytes','KiB','MiB','GiB']


fileUploadInput.addEventListener("change",async (event) => {
    event.preventDefault();
    const fileUploaded = fileUploadInput.files[0]
    let fileSize = fileUploaded.size
    let sizeUnit = 0
    while (fileSize > 1024) {
        sizeUnit += 1
        fileSize /= 1024
    }
    fileSize = Math.round((fileSize + Number.EPSILON) * 100) / 100
    document.getElementById("selectedFile").innerHTML = `<p>Uploaded file : ${fileUploaded.name} with size ${fileSize} ${sizeUnits[sizeUnit]}</p>`
});

fileForm.addEventListener("submit", async (event) => {
    event.preventDefault()
    const formData = new FormData(fileForm)

    const response = await fetch("/upload", {
        method: "POST",
        body: formData
    });

    if (response.ok) {
        const result = await response.text();
        document.getElementById("downloadFile").innerHTML = result
    } else {
        document.getElementById("downloadFile").innerHTML = "<p>Error processing File!</p>"
    }
})

