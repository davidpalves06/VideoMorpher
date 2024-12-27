import {FileUploadScript} from "./FileUpload.mjs"

const pageHeader = document.getElementById("VideoMorpherHeader")

pageHeader.addEventListener('click',(event) => {
    location.reload()
})
FileUploadScript()




