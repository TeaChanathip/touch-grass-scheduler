"use client"

import Image from "next/image"
import MenuRoundedIcon from "@mui/icons-material/MenuRounded"
import Link from "next/link"
import React, { useEffect, useState } from "react"
import { usePathname, useRouter } from "next/navigation"
import { useAppDispatch, useAppSelector } from "../store/hooks"
import {
    userGetMe,
    userLogout,
    selectUser,
    selectUserStatus,
} from "../store/features/user/userSlice"
import { User, UserRole } from "../interfaces/User.interface"
import MyButton from "../components/MyButton"

interface Path {
    title: string
    path: string
    visiblility: "all" | "unauthenticated" | "authenticated" | UserRole
}

interface NavItem extends Path {
    id: string
}

export default function Navbar() {
    // Store
    const dispatch = useAppDispatch()
    const userStatus = useAppSelector(selectUserStatus)
    const user = useAppSelector(selectUser)

    // Other hooks
    const [isShow, setShow] = useState(false)
    const pathname = usePathname()
    const router = useRouter()

    useEffect(() => {
        setShow(false)
    }, [pathname])

    // TODO: Migrate to use Reach Server Component?
    useEffect(() => {
        dispatch(userGetMe())
    }, [dispatch])

    useEffect(() => {
        if (userStatus === "unauthenticated") {
            router.push("/login")
        }
    }, [userStatus, router])

    const routes: Path[] = [
        { title: "Login", path: "/login", visiblility: "unauthenticated" },
        { title: "RouteAll", path: "/all", visiblility: "all" },
        {
            title: "RouteUnauthenticated",
            path: "/unauthenticated",
            visiblility: "unauthenticated",
        },
        {
            title: "RouteAuthenticated",
            path: "/authenticated",
            visiblility: "authenticated",
        },
        {
            title: "RouteStudent",
            path: "/student",
            visiblility: UserRole.STUDENT,
        },
        {
            title: "RouteTeacher",
            path: "/teacher",
            visiblility: UserRole.TEACHER,
        },
        {
            title: "RouteGuardian",
            path: "/guardian",
            visiblility: UserRole.GUARDIAN,
        },
        { title: "RouteAdmin", path: "/admin", visiblility: UserRole.ADMIN },
    ]

    // Generate unique id to be used as key
    const routeItems: NavItem[] = routes
        .map((route) => {
            return { ...route, id: `route-items-${route.title}` }
        })
        .filter((item) => {
            if (item.path === pathname) return false

            switch (item.visiblility) {
                case "all":
                    return true
                case "unauthenticated":
                    return userStatus != "authenticated"
                case "authenticated":
                    return userStatus === "authenticated"
                default:
                    if (userStatus !== "authenticated") return false
                    return user?.role === item.visiblility
            }
        })

    return (
        <>
            <header className="sticky top-0 h-14 bg-prim-green-800 flex items-center justify-between px-2.5">
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
                <button
                    onClick={() => {
                        setShow(!isShow)
                    }}
                >
                    <MenuRoundedIcon
                        fontSize="large"
                        className="text-white [&:hover,&:focus]:text-prim-green-500 cursor-pointer"
                    />
                </button>
            </header>
            <NavPanel navItems={routeItems} isShow={isShow} user={user} />
        </>
    )
}

function NavPanel({
    navItems,
    isShow,
    user,
}: {
    navItems: NavItem[]
    isShow: boolean
    user?: User
}) {
    // Store
    const dispatch = useAppDispatch()

    // Other hooks
    const router = useRouter()

    let displayName = ""
    if (user) {
        displayName = user.first_name
        if (user.middle_name) displayName += " " + user.middle_name
        if (user.last_name) displayName += " " + user.last_name

        const { innerWidth } = window
        if (innerWidth < 768 && displayName.length > 16) {
            displayName = displayName.slice(0, 16).trim() + "..."
        } else if (innerWidth >= 768 && displayName.length > 40) {
            displayName = displayName.slice(0, 40).trim() + "..."
        }
    }

    // Button Handlers
    const editBtnHandler = () => {
        router.push("/profile")
    }

    const logoutBtnHandler = () => {
        dispatch(userLogout())
    }

    return (
        <aside
            className="fixed h-full w-full lg:w-[500px] right-0 bg-prim-green-500/90 transition-all ease-in overflow-y-auto custom-scrollbar"
            style={
                isShow
                    ? { transform: "none" }
                    : { transform: "translateX(100%)" }
            }
        >
            {user && (
                <div className="mt-8 mx-8 px-4 py-3 flex flex-row items-center bg-prim-green-100 rounded-xl text-2xl drop-shadow-md">
                    <Image
                        src={user.avatar_url ?? "default_avartar.svg"}
                        alt="avartar"
                        width={120}
                        height={120}
                        className="size-[120px] rounded-full"
                    />
                    <div className="w-full flex flex-col gap-6">
                        <p className="w-full text-center text-xl">
                            {displayName}
                        </p>
                        <span className="w-full flex flex-wrap justify-center gap-2">
                            <MyButton
                                variant="positive"
                                className="text-sm md:text-xl w-20 md:w-28"
                                onClick={() => {
                                    editBtnHandler()
                                }}
                            >
                                Edit
                            </MyButton>
                            <MyButton
                                variant="negative"
                                className="text-sm md:text-xl w-20 md:w-28"
                                onClick={() => {
                                    logoutBtnHandler()
                                }}
                            >
                                Logout
                            </MyButton>
                        </span>
                    </div>
                </div>
            )}
            <nav className="mt-4">
                <ul className="mx-8 flex flex-col gap-4">
                    {navItems.map((item) => (
                        <li
                            key={item.id}
                            className="bg-prim-green-100 [&:hover,&:focus]:bg-prim-green-500 h-14 rounded-xl text-2xl drop-shadow-md"
                        >
                            <Link
                                href={item.path}
                                className="size-full px-5 flex items-center"
                            >
                                {item.title}
                            </Link>
                        </li>
                    ))}
                </ul>
            </nav>
        </aside>
    )
}

// TODO: Fix tabIndex of NavItem
