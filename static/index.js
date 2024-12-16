const fileForm = document.getElementById("fileForm");
const fileUploadInput = document.getElementById("fileUpload");
const progressTracker = document.getElementById("progressTracker");
const downloadFile = document.getElementById("downloadFile");
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
        const resultJSON = await response.json();
        const processID = resultJSON.processID
        trackVideoProgress(processID)
        downloadFile.innerHTML = resultJSON.downloadRef
        downloadFile.hidden = true
    } else {
        document.getElementById("downloadFile").innerHTML = "<p>Error processing File!</p>"
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

