import React from "react"
import GuestRoute from "../../../layout/GuestRoute"
import PageLayout from "../../../layout/PageLayout"

export default function ResetPwdLayout({
    children,
}: {
    children: React.ReactNode
}) {
    return (
        <GuestRoute>
            <PageLayout title="Reset Password">{children}</PageLayout>
        </GuestRoute>
    )
}
