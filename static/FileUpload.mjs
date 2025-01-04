import { ffmpegSupportedFormats, LOCAL_STORAGE_DOWNLOADS_KEY, sizeUnits } from "./Constants.mjs";
import { OperationInputScript } from "./OperationInput.mjs"

const fileForm = document.getElementById("fileForm");
const fileUploadInput = document.getElementById("fileUpload");
const formButton = document.getElementById("formButton")
const downloadFile = document.getElementById("downloadFile");
const SelectedFileInfo = document.getElementById("selectedFile");
const progressContainer = document.getElementById("progressContainer");
const progressTracker = document.getElementById("progressTracker");
const inputVideoPlayer = document.getElementById("inputVideoPlayer");
const operationSelectionContainer = document.getElementById("operationSelectionContainer")
const outputVideoPlayer = document.getElementById("outputVideoPlayer")
const progressLabel = document.getElementById('progressLabel')
let streamLink = ''
let outputFileName = ''


const handleFileUpload = () => {
    let videoErrorMessage = document.getElementById("inputVideoErrorMessage")
    videoErrorMessage.innerText = ""
    if (fileUploadInput.files && fileUploadInput.files[0]) {
        const fileUploaded = fileUploadInput.files[0]
        let fileName = fileUploaded.name
        let fileExtension = fileName.split(".").pop().toLowerCase()
        if (!ffmpegSupportedFormats.includes(fileExtension)) {
            fileUploadInput.value = ''
            SelectedFileInfo.innerHTML = "Can't upload this file format"
            videoErrorMessage.innerText = "File format not accepted"
            videoErrorMessage.hidden = false
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
        inputVideoPlayer.innerHTML = '';
        let videoSource = document.createElement('source')
        videoSource.src = URL.createObjectURL(fileUploaded)
        videoSource.addEventListener('error', (event) => {
            inputVideoPlayer.hidden = true
            videoErrorMessage.hidden = false
            videoErrorMessage.innerText = "File format cannot be played"
        })

        inputVideoPlayer.appendChild(videoSource)
        inputVideoPlayer.hidden = false
        inputVideoPlayer.load()
        videoErrorMessage.hidden = true
        operationSelectionContainer.hidden = false
        OperationInputScript()
        operationSelectionContainer.style.display = 'flex'
        formButton.hidden = false
    }
    else {
        inputVideoPlayer.hidden = true
        operationSelectionContainer.hidden = true
        operationSelectionContainer.style.display = "none"
        formButton.hidden = true
    }
}


const addToLocalStorageList = () => {
    let pastDownloads = localStorage.getItem(LOCAL_STORAGE_DOWNLOADS_KEY);
    let newFile = {
        outputFileName,
        creationDate: new Date()
    }
    if (pastDownloads == null) {
        let pastDownloadList = {
            downloads: []
        }

        pastDownloadList.downloads.push(newFile)

        localStorage.setItem(LOCAL_STORAGE_DOWNLOADS_KEY, JSON.stringify(pastDownloadList))
    } else {
        let pastDownloadList = JSON.parse(pastDownloads)
        pastDownloadList.downloads.push(newFile)
        localStorage.setItem(LOCAL_STORAGE_DOWNLOADS_KEY, JSON.stringify(pastDownloadList))
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
        const progressBeginning = document.getElementById("progressBeginning")
        fileForm.style.display = "none"
        progressLabel.innerHTML = `0%`
        progressContainer.style.display = "block"
        progressBeginning.hidden = false
        const response = await fetch("/upload", {
            method: "POST",
            body: formData
        });

        if (response.ok) {
            const resultJSON = await response.json();
            const processID = resultJSON.processID
            trackVideoProgress(processID)
            downloadFile.innerHTML = `<a href='/download?file=${resultJSON.generatedFile}&stream=disabled' class="chooseFileLabel">Download file</a>`
            downloadFile.hidden = true
            progressBeginning.hidden = true
            outputFileName = resultJSON.generatedFile
            streamLink = `/download?file=${outputFileName}&stream=enabled`
        } else {
            downloadFile.innerHTML = "<p>Error processing File!</p>"
        }
    })

    function trackVideoProgress(processID) {
        const progressSource = new EventSource(`/progress?processID=${processID}`)


        progressSource.addEventListener('progress', function (event) {
            const progress = event.data

            if (progress < 100) {
                progressLabel.innerHTML = `<em>${progress}%</em>`
                progressLabel.style.left = `${progress - 2}%`
                progressTracker.style.width = `${progress}%`
            } else {
                progressLabel.innerHTML = `<em>100%</em>`
                progressLabel.style.left = `98%`
                progressTracker.style.width = `100%`
                progressSource.close()
                document.getElementById("videoProcessedInfo").hidden = false
                let videoSource = document.createElement('source')
                videoSource.src = streamLink
                videoSource.onerror = ((event) => {
                    outputVideoPlayer.innerHTML = ""
                    outputVideoPlayer.hidden = true
                    let outputVideoErrorMessage = document.getElementById("outputVideoErrorMessage")
                    outputVideoErrorMessage.innerText = "File format cannot be played"
                    outputVideoErrorMessage.hidden = false
                })
                downloadFile.hidden = false
                outputVideoPlayer.appendChild(videoSource)
                outputVideoPlayer.hidden = false
                addToLocalStorageList(outputFileName)
            }
        });



        progressSource.onerror = (event) => {
            progressSource.close()
            let outputVideoErrorMessage = document.getElementById("outputVideoErrorMessage")
            outputVideoErrorMessage.innerHTML = "<strong>Error while processing video. Try again please</strong>"
            outputVideoErrorMessage.hidden = false
            progressContainer.style.display = 'none'
        }
    }
}


fileUploadInput.value = ""
fileUploadInput.files = undefined

handleFileUpload()