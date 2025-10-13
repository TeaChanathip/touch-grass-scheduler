import React from "react"
import "../styles/globals.css"
import Navbar from "../layout/Navbar"

export default function RootLayout({
    children,
}: {
    children: React.ReactNode
}) {
    return (
        <>
            <Navbar />
            <html lang="en">
                <body>{children}</body>
            </html>
        </>
    )
}
