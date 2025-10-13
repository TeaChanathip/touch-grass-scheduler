"use client"

import Image from "next/image"
import MenuRoundedIcon from "@mui/icons-material/MenuRounded"
import Link from "next/link"
import React, { useEffect, useState } from "react"
import { usePathname } from "next/navigation"

interface Path {
    title: string
    path: string
}

interface NavItem extends Path {
    id: string
}

export default function Navbar() {
    const [isShow, setShow] = useState(false)
    const pathname = usePathname()

    useEffect(() => {
        setShow(false)
    }, [pathname])

    const routes: Path[] = [
        { title: "Login", path: "/login" },
        { title: "Route1", path: "/" },
        { title: "Route2", path: "/" },
        { title: "Route3", path: "/" },
        { title: "Route4", path: "/" },
        { title: "Route5", path: "/" },
        { title: "Route6", path: "/" },
    ]

    // Generate unique id to be used as key
    const routeItems: NavItem[] = routes.map((route) => {
        return { ...route, id: crypto.randomUUID() }
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
            <NavPanel navItems={routeItems} isShow={isShow} />
        </>
    )
}

function NavPanel(props: { navItems: NavItem[]; isShow: boolean }) {
    const pathname = usePathname()

    return (
        <aside
            className="fixed h-full w-full lg:w-[400px] right-0 bg-prim-green-500/90 transition-all ease-in overflow-y-auto custom-scrollbar"
            style={
                props.isShow
                    ? { transform: "none" }
                    : { transform: "translateX(100%)" }
            }
        >
            <nav>
                <ul className="mt-8 mx-8 flex flex-col gap-4">
                    {props.navItems
                        .filter((item) => item.path != pathname)
                        .map((item) => (
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

// TODO: Hide Login route when User is already logged in
// TODO: Show user info component when logged in
