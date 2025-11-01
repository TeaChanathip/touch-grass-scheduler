function convertToWebP(file: File, quality: number = 0.9): Promise<File> {
    return new Promise((resolve, reject) => {
        // 1. Create an Image object
        const img = new Image()

        // 2. Load the original file into the Image object
        img.src = URL.createObjectURL(file)

        img.onload = () => {
            // 3. Create a hidden canvas
            const canvas = document.createElement("canvas")
            canvas.width = img.width
            canvas.height = img.height

            // 4. Draw the image onto the canvas
            const ctx = canvas.getContext("2d")
            if (!ctx) {
                return reject(new Error("Failed to get canvas context"))
            }
            ctx.drawImage(img, 0, 0)

            // 5. Export the canvas content as a WebP blob
            canvas.toBlob(
                (blob) => {
                    if (!blob) {
                        return reject(new Error("Canvas toBlob failed"))
                    }

                    // 6. Create a new File object from the blob
                    const webpFile = new File([blob], "converted.webp", {
                        type: "image/webp",
                    })

                    resolve(webpFile)
                },
                "image/webp", // Specify the output format
                quality // Specify the quality
            )
        }

        img.onerror = (err) => {
            console.log(err)
            reject("convert to WebP failed")
        }
    })
}

export default convertToWebP
