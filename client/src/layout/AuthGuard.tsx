"use client"

import React, { useEffect } from "react"
import { UserRole } from "../interfaces/User.interface"
import { useAppSelector } from "../store/hooks"
import { selectUser, selectUserStatus } from "../store/features/user/userSlice"
import { useRouter } from "next/navigation"
import { CircularProgress } from "@mui/material"

export default function AuthGuard({
    roles,
    children,
}: {
    roles?: UserRole[]
    children: React.ReactNode
}) {
    // Store
    const userStatus = useAppSelector(selectUserStatus)
    const user = useAppSelector(selectUser)

    // Hooks
    const router = useRouter()

    useEffect(() => {
        const isUnauthenticated: boolean =
            userStatus === "unauthenticated" || userStatus === "error"
        const isUnauthorized: boolean =
            roles !== undefined &&
            userStatus === "authenticated" &&
            user !== undefined &&
            !roles.includes(user.role)

        if (isUnauthenticated || isUnauthorized) {
            router.replace("/login")
        }
    }, [userStatus, user, roles, router])

    if (userStatus === "idle" || userStatus === "loading") {
        return (
            <span className="text-prim-green-400 mt-10">
                <CircularProgress color="inherit" />
            </span>
        )
    }

    if (userStatus === "unauthenticated" || userStatus === "error") {
        return null
    }

    return <>{children}</>
}
