import { LOCAL_STORAGE_DOWNLOADS_KEY } from "./Constants.mjs"
import { FileUploadScript } from "./FileUpload.mjs"

const pageHeader = document.getElementById("VideoMorpherHeader")
const pastDownloads = document.getElementById("pastDownloadList")

pageHeader.addEventListener('click', (event) => {
    location.reload()
})

function isMoreThanFiveMinutesBefore(date1, date2) {
    let diffMs = date1.getTime() - date2.getTime();

    let diffMins = Math.abs(diffMs / 1000 / 60);

    return diffMins > 5;
}

let storageDownloads = localStorage.getItem(LOCAL_STORAGE_DOWNLOADS_KEY);

if (storageDownloads != null) {
    let pastDownloadList = JSON.parse(storageDownloads)
    let filteredPastDownloads = pastDownloadList.downloads.filter((file) => !isMoreThanFiveMinutesBefore(new Date(file.creationDate), new Date()));

    for (const file of filteredPastDownloads) {
        let newListItem = document.createElement('li')
        newListItem.innerHTML = `<a href="/download?file=${file.outputFileName}&stream=disabled">${file.outputFileName}</a>`
        pastDownloads.appendChild(newListItem)
    }

    localStorage.setItem(LOCAL_STORAGE_DOWNLOADS_KEY,JSON.stringify({downloads:filteredPastDownloads}))
}

FileUploadScript()




