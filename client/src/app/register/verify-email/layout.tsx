import React from "react"
import PageLayout from "../../../layout/PageLayout"
import GuestRoute from "../../../layout/GuestRoute"

export default function LoginLayout({
    children,
}: {
    children: React.ReactNode
}) {
    return (
        <GuestRoute>
            <PageLayout title="Verify Email">{children}</PageLayout>
        </GuestRoute>
    )
}
