"use client"

import Image from "next/image"
import MenuRoundedIcon from "@mui/icons-material/MenuRounded"
import Link from "next/link"
import React, { memo, useEffect, useState } from "react"
import { usePathname, useRouter } from "next/navigation"
import { useAppDispatch, useAppSelector } from "../store/hooks"
import {
    userAutoLogin,
    userLogout,
    selectUser,
    selectUserStatus,
} from "../store/features/user/userSlice"
import { UserRole } from "../interfaces/User.interface"
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

    // Hooks
    const pathname = usePathname()
    const [isSidebarShown, setSidebarShown] = useState(false)

    useEffect(() => {
        dispatch(userAutoLogin())
    }, [dispatch])

    useEffect(() => {
        setSidebarShown(false)
    }, [pathname])

    return (
        <>
            <header className="sticky top-0 h-14 bg-prim-green-800 flex items-center justify-between px-2.5">
                <AppIcon />
                <button
                    onClick={() => {
                        setSidebarShown(!isSidebarShown)
                    }}
                >
                    <MenuRoundedIcon
                        fontSize="large"
                        className="text-white [&:hover,&:focus]:text-prim-green-500 cursor-pointer"
                    />
                </button>
            </header>
            <Sidebar isSidebarShown={isSidebarShown} />
        </>
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

const Sidebar = memo(function Sidebar({
    isSidebarShown,
}: {
    isSidebarShown: boolean
}) {
    return (
        <aside
            inert={!isSidebarShown}
            className="fixed h-full w-full lg:w-[500px] right-0 bg-prim-green-500/90 transition-all ease-in overflow-y-auto custom-scrollbar"
            style={
                isSidebarShown
                    ? { transform: "none" }
                    : { transform: "translateX(100%)" }
            }
        >
            <UserCard />
            <NavMenu />
        </aside>
    )
})

const UserCard = memo(function UserCard() {
    // Store
    const dispatch = useAppDispatch()
    const user = useAppSelector(selectUser)

    // Hooks
    const router = useRouter()

    // Button Handlers
    const editBtnHandler = () => {
        router.push("/profile")
    }
    const logoutBtnHandler = () => {
        dispatch(userLogout())
    }

    if (user === undefined) {
        return null
    }

    return (
        <div className="mt-8 mx-8 px-4 py-3 flex flex-row items-center bg-prim-green-100 rounded-xl text-2xl drop-shadow-md">
            <Image
                src={user.avatar_url ?? "default_avatar.svg"}
                alt="avatar"
                width={120}
                height={120}
                className="size-[120px] rounded-full"
            />
            <div className="w-full flex flex-col gap-6 overflow-hidden">
                <p className="w-full text-center text-xl whitespace-nowrap overflow-hidden text-ellipsis">
                    {user.first_name} {user.middle_name} {user.last_name}
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
    )
})

// All available routes
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

const NavMenu = memo(function NavMenu() {
    // Hooks
    const pathname = usePathname()
    const userStatus = useAppSelector(selectUserStatus)
    const user = useAppSelector(selectUser)

    // Generate unique id to be used as key
    // Filter routes based on the role of state of user
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
        <nav className="mt-4">
            <ul className="mx-8 flex flex-col gap-4">
                {routeItems.map((item) => (
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
    )
})
