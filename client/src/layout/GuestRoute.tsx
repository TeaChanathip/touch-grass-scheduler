"use client"

import React, { useEffect } from "react"
import { useAppSelector } from "../store/hooks"
import { selectUserStatus } from "../store/features/user/userSlice"
import { useRouter } from "next/navigation"

export default function GuestRoute({
    children,
}: {
    children: React.ReactNode
}) {
    // Store
    const userStatus = useAppSelector(selectUserStatus)

    // Hooks
    const router = useRouter()

    useEffect(() => {
        if (userStatus === "authenticated") {
            router.replace("/")
        }
    }, [userStatus, router])

    return children
}
