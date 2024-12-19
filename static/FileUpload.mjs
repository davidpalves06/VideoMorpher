import { ffmpegSupportedFormats, sizeUnits } from "./Constants.mjs";

const fileForm = document.getElementById("fileForm");
const fileUploadInput = document.getElementById("fileUpload");
const downloadFile = document.getElementById("downloadFile");
const SelectedFileInfo = document.getElementById("selectedFile");
const progressTracker = document.getElementById("progressTracker");
const inputVideoPlayer = document.getElementById("inputVideoPlayer");

const handleFileUpload = () => {
    if (fileUploadInput.files && fileUploadInput.files[0]) {
        const fileUploaded = fileUploadInput.files[0]
        let fileName = fileUploaded.name
        let fileExtension = fileName.split(".").pop().toLowerCase()
        if (!ffmpegSupportedFormats.includes(fileExtension)) {
            fileUploadInput.value = ''
            SelectedFileInfo.innerHTML = "Can't upload this file format"
            return
        }
        let fileSize = fileUploaded.size
        let sizeUnit = 0
        while (fileSize > 1024) {
            sizeUnit += 1
            fileSize /= 1024
        }
        fileSize = Math.round((fileSize + Number.EPSILON) * 100) / 100
        document.getElementById("selectedFile").innerHTML = `<strong><em>${fileUploaded.name}</em></strong> with <strong>${fileSize}</strong> <em>${sizeUnits[sizeUnit]}</em>`
        document.getElementById("chooseFileLabel").innerText = "Change File";
        inputVideoPlayer.src = URL.createObjectURL(fileUploaded)
        inputVideoPlayer.hidden = false
    }
}

export const FileUploadScript = () => {
    fileUploadInput.accept = ffmpegSupportedFormats.map((format) => '.' + format).join(',')

    fileUploadInput.addEventListener("change", (event) => {
        event.preventDefault();
        handleFileUpload()
    });

    fileForm.addEventListener("submit", async (event) => {
        event.preventDefault()
        const formData = new FormData(fileForm)

        const response = await fetch("/upload", {
            method: "POST",
            body: formData
        });

        if (response.ok) {
            const resultJSON = await response.json();
            const processID = resultJSON.processID
            trackVideoProgress(processID)
            downloadFile.innerHTML = resultJSON.downloadRef
            downloadFile.hidden = true
        } else {
            downloadFile.innerHTML = "<p>Error processing File!</p>"
        }
    })

    function trackVideoProgress(processID) {
        const progressSource = new EventSource(`/progress?processID=${processID}`)


        progressSource.onopen = (event) => {
            progressTracker.innerHTML = `<p>Progress: 0%</p>`
        }

        progressSource.addEventListener('progress', function (event) {
            const progress = event.data

            if (progress < 100) {
                progressTracker.innerHTML = `<p>Progress: ${progress} %</p>`
            } else {
                progressSource.close()
                progressTracker.innerHTML = `<p>Video processed! You can download it in the link below!</p>`
                downloadFile.hidden = false
            }
        });



        progressSource.onerror = (event) => {
            progressSource.close()
        }
    }
}

handleFileUpload()