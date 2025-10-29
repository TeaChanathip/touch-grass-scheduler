import React from "react"
import PageLayout from "../../layout/PageLayout"
import AuthGuard from "../../layout/AuthGuard"
import { UserRole } from "../../interfaces/User.interface"

export default function LoginLayout({
    children,
}: {
    children: React.ReactNode
}) {
    return (
        <AuthGuard
            roles={[UserRole.STUDENT, UserRole.TEACHER, UserRole.GUARDIAN]}
        >
            <PageLayout title="Profile">{children}</PageLayout>
        </AuthGuard>
    )
}
