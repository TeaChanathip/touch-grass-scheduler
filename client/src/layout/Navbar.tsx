"use client"

import Image from "next/image"
import Link from "next/link"
import React, { memo, useEffect } from "react"
import { useAppDispatch, useAppSelector } from "../store/hooks"
import {
    selectUserStatus,
    userAutoLogin,
} from "../store/features/user/userSlice"
import NotiButton from "./NotiButton"
import SidebarButton from "./SidebarButton"

export default function Navbar() {
    // Store
    const dispatch = useAppDispatch()
    const userStatus = useAppSelector(selectUserStatus)

    useEffect(() => {
        dispatch(userAutoLogin())
    }, [dispatch])

    return (
        <div>
            <header className="sticky top-0 h-14 bg-prim-green-800 flex items-center justify-between px-2.5">
                <AppIcon />
                <span className="h-full flex items-center gap-6">
                    {userStatus === "authenticated" && <NotiButton />}
                    <SidebarButton />
                </span>
            </header>
        </div>
    )
}

const AppIcon = memo(function AppIcon() {
    return (
        <Link href="/">
            <Image
                src="icon.svg"
                alt="icon"
                width={40}
                height={40}
                className="inline mr-2"
            />
            <Image
                src="text_icon_light.svg"
                alt="text-icon"
                width={90}
                height={40}
                className="inline size-auto"
            />
        </Link>
    )
})
