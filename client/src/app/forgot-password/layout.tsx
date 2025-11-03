import React from "react"
import PageLayout from "../../layout/PageLayout"
import GuestRoute from "../../layout/GuestRoute"

export default function ForgotPasswordLayout({
    children,
}: {
    children: React.ReactNode
}) {
    return (
        <GuestRoute>
            <PageLayout title="Forgot Password">{children}</PageLayout>
        </GuestRoute>
    )
}
