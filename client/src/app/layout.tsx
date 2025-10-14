import React from "react"
import "../styles/globals.css"
import Navbar from "../layout/Navbar"
import StoreProvider from "./StoreProvider"

export default function RootLayout({
    children,
}: {
    children: React.ReactNode
}) {
    return (
        <StoreProvider>
            <Navbar />
            <html lang="en">
                <body>{children}</body>
            </html>
        </StoreProvider>
    )
}
