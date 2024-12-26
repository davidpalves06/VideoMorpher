import { ffmpegSupportedFormats } from "./Constants.mjs";

function updateRangeLabel(params) {
    const range = document.getElementById("speedSlider")
    const rangeLabel = document.getElementById("speedLabel")
    const newValue = Number(((range.value - range.min) * 100) / (range.max - range.min));
    const newPosition = 10 + (newValue * 0.1);
    
    rangeLabel.innerHTML = Number(range.value).toFixed(2);
    rangeLabel.style.left = `calc(${newValue}% - ${newPosition}px)`;
}

export const OperationInputScript = () => {
    const fileForm = document.getElementById("fileForm");
    const operationSelectors = fileForm.querySelectorAll('input[name="operation"]');
    const operationInput = document.getElementById("operationInput");

    function changeOperationInput() {
        const operationSelector = fileForm.querySelector('input[name="operation"]:checked');
        operationInput.innerHTML = ""
        if (operationSelector.value === "conversion") {
            const select = document.createElement('select');
            //TODO: REMOVE FOR FROM FILE SELECTED
            select.id = 'conversionFormat';
            select.name = 'conversionFormat';
            select.className = 'formatSelect'

            ffmpegSupportedFormats.forEach(format => {
                const option = document.createElement('option');
                option.value = format;
                option.textContent = format;
                select.appendChild(option);
            });

            operationInput.appendChild(select);
        } else if (operationSelector.value === "motion") {
            const speedSlider = document.createElement('input')
            speedSlider.id = "speedSlider"
            speedSlider.className = "motionSliderInput"
            speedSlider.type = 'range'
            speedSlider.min = 0.25
            speedSlider.max = 10
            speedSlider.step = 0.01
            speedSlider.value = Number(1).toFixed(2)
            speedSlider.name = "motionSpeed"
            const speedIndicator = document.createElement('label')
            speedIndicator.className = "speedRangeIndicator"
            speedIndicator.id = "speedLabel"

            

            speedSlider.addEventListener('input', updateRangeLabel)
            operationInput.appendChild(speedSlider)
            operationInput.appendChild(speedIndicator)
            updateRangeLabel()

        } else if (operationSelector.value === "reverse") {
            operationInput.innerHTML = "REVERSE"
        }
    }

    operationSelectors.forEach((radio) => radio.addEventListener('change', (event) => {
        changeOperationInput()
    }))

    changeOperationInput()
}
