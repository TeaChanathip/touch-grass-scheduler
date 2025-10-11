"use client"

import Image from "next/image"
import MenuRoundedIcon from "@mui/icons-material/MenuRounded"
import Link from "next/link"
import { useEffect, useState } from "react"
import { usePathname } from "next/navigation"

interface Route {
    title: string
    route: string
}

interface RouteItem extends Route {
    id: string
}

export default function Navbar() {
    const [isShow, setShow] = useState(false)
    const pathname = usePathname()

    useEffect(() => {
        setShow(false)
    }, [pathname])

    const routes: Route[] = [
        { title: "Login", route: "/login" },
        { title: "Route1", route: "/" },
        { title: "Route2", route: "/" },
        { title: "Route3", route: "/" },
    ]

    // Generate unique id to be used as key
    const routeItems: RouteItem[] = routes.map((route) => {
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
                        className="inline"
                    />
                </Link>
                <button
                    onClick={() => {
                        setShow(!isShow)
                    }}
                >
                    <MenuRoundedIcon fontSize="large" className="text-white" />
                </button>
            </header>
            <NavPanel routeItems={routeItems} isShow={isShow} />
        </>
    )
}

function NavPanel(props: { routeItems: RouteItem[]; isShow: boolean }) {
    return (
        <aside
            className="fixed h-full w-full bg-prim-green-500/90 transition-all ease-in"
            style={{ left: props.isShow ? "0" : "100vw" }}
        >
            <nav>
                <ul className="mt-8 mx-8 flex flex-col gap-4 text-prim-dark">
                    {props.routeItems.map((routeItem) => (
                        <li
                            key={routeItem.id}
                            className="bg-prim-green-100 h-14 rounded-xl text-2xl"
                        >
                            <Link
                                href={routeItem.route}
                                className="size-full px-5 flex items-center"
                            >
                                {routeItem.title}
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
