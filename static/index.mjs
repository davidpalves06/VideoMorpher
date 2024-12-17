import { ffmpegSupportedFormats,sizeUnits } from "./Constants.mjs";

const fileForm = document.getElementById("fileForm");
const fileUploadInput = document.getElementById("fileUpload");
const progressTracker = document.getElementById("progressTracker");
const downloadFile = document.getElementById("downloadFile");
const operationSelector = document.getElementById("operationSelect");
const operationInput = document.getElementById("operationInput");
const SelectedFileInfo = document.getElementById("selectedFile");

fileUploadInput.accept = ffmpegSupportedFormats.map((format) => '.'+format).join(',')

fileUploadInput.addEventListener("change",async (event) => {
    event.preventDefault();
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
    document.getElementById("selectedFile").innerHTML = `<p>Uploaded file : ${fileUploaded.name} with size ${fileSize} ${sizeUnits[sizeUnit]}</p>`
});

fileForm.addEventListener("submit", async (event) => {
    event.preventDefault()
    const formData = new FormData(fileForm)
    console.log(formData);
    
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

    progressSource.addEventListener('progress',function(event) {
        const progress = event.data
        if (progress !== "100%") {
            progressTracker.innerHTML = `<p>Progress: ${progress}</p>`
        } else {
            progressSource.close()
            progressTracker.innerHTML = `<p>Video processed! You can download it in the link below!</p>`
            downloadFile.hidden = false
        }
    });



    progressSource.onerror = (event) => {
        console.log('EventSource connection state:', progressSource.readyState);
        progressSource.close()
    }   
}

function changeOperationInput() {
    operationInput.innerHTML = ""
    if (operationSelector.value === "conversion") {
        const select = document.createElement('select');
        select.id = 'conversionFormat';

        ffmpegSupportedFormats.forEach(format => {
            const option = document.createElement('option');
            option.value = format;
            option.textContent = format;
            select.appendChild(option);
        });

        operationInput.appendChild(select);
    } else if (operationSelector.value === "motion") {
        const speedSlider = document.createElement('input')
        speedSlider.type = 'range'
        speedSlider.min = 0.25
        speedSlider.value = 1
        speedSlider.max = 2
        speedSlider.step = 0.01
        speedSlider.name = "motionSpeed"
        const speedIndicator = document.createElement('p')
        speedIndicator.textContent = 'Speed value : 1'
        
        speedSlider.addEventListener('change',(event) => {
            speedIndicator.textContent = 'Speed value :' + speedSlider.value
        })
        operationInput.appendChild(speedSlider)
        operationInput.appendChild(speedIndicator)

    } else if (operationSelector.value === "reverse") {
        operationInput.innerHTML = "REVERSE"
    }
}

operationSelector.addEventListener('change',(event) => {
    changeOperationInput()
})

changeOperationInput()

