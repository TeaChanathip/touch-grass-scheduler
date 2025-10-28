import React from "react"
import PageLayout from "../../layout/PageLayout"

export default function LoginLayout({
    children,
}: {
    children: React.ReactNode
}) {
    return <PageLayout title="Login">{children}</PageLayout>
}
