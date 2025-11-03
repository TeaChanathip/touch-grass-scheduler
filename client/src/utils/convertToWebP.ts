interface Options {
    quality?: number
    square?: boolean
    maxWidth?: number // Used when sqaure is true
}

function convertToWebP(
    file: File,
    { quality = 0.9, square, maxWidth }: Options = {}
): Promise<File> {
    return new Promise((resolve, reject) => {
        // 1. Create an Image object
        const img = new Image()

        // 2. Load the original file into the Image object
        img.src = URL.createObjectURL(file)

        img.onload = () => {
            // Source
            let sx = 0
            let sy = 0
            let sWidth = img.width
            let sHeight = img.height

            // Destination
            let dWidth: number
            let dHeight: number

            if (square) {
                // 1. Find the largest possible square from the original.
                sWidth = Math.min(img.width, img.height)
                sHeight = sWidth

                // 2. Find the (x, y) offset to center this square.
                sx = Math.trunc((img.width - sWidth) / 2)
                sy = Math.trunc((img.height - sHeight) / 2)

                // 3. Set the *final* canvas size.
                dWidth = maxWidth || sWidth
                dHeight = maxWidth || sHeight
            } else {
                // 4. No changes, use original dimensions
                dWidth = img.width
                dHeight = img.height
            }

            // 5. Create a hidden canvas
            const canvas = document.createElement("canvas")
            canvas.width = dWidth
            canvas.height = dHeight

            // 6. Draw the image onto the canvas
            const ctx = canvas.getContext("2d")
            if (!ctx) {
                return reject(new Error("Failed to get canvas context"))
            }
            ctx.drawImage(img, sx, sy, sWidth, sHeight, 0, 0, dWidth, dHeight)

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
                "image/webp",
                quality
            )
        }

        img.onerror = (err) => {
            console.log(err)
            reject("convert to WebP failed")
        }
    })
}

export default convertToWebP
