import { memo, useEffect, useRef, useState } from "react"
import { useAppDispatch, useAppSelector } from "../store/hooks"
import {
    selectUser,
    selectUserStatus,
    userLogout,
} from "../store/features/user/userSlice"
import { usePathname, useRouter } from "next/navigation"
import Image from "next/image"
import MyButton from "../components/MyButton"
import { UserRole } from "../interfaces/User.interface"
import Link from "next/link"
import useClickOutside from "../hooks/useClickOutside"

import MenuRoundedIcon from "@mui/icons-material/MenuRounded"

const SidebarButton = memo(function SidebarButton() {
    // Hooks
    const [isSidebarShown, setSidebarShown] = useState(false)
    const sidebarRef = useRef<HTMLDivElement | null>(null)
    useClickOutside(sidebarRef, () => setSidebarShown(false))
    const pathname = usePathname()

    useEffect(() => {
        setSidebarShown(false)
    }, [pathname])

    return (
        <div ref={sidebarRef}>
            <button
                onClick={() =>
                    setSidebarShown((isSidebarShown) => !isSidebarShown)
                }
                className="text-white hover:text-prim-green-500 active:text-prim-green-500 cursor-pointer"
            >
                <MenuRoundedIcon fontSize="large" />
            </button>
            <Sidebar isSidebarShown={isSidebarShown} />
        </div>
    )
})
export default SidebarButton

function Sidebar({ isSidebarShown }: { isSidebarShown: boolean }) {
    return (
        <aside
            inert={!isSidebarShown}
            className="fixed h-full w-full lg:w-[500px] right-0 bg-prim-green-500/90 
                transition-all ease-in overflow-y-auto custom-scrollbar mt-[10px]"
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
}

const UserCard = memo(function UserCard() {
    // Store
    const dispatch = useAppDispatch()
    const user = useAppSelector(selectUser)

    // Hooks
    const router = useRouter()
    const pathname = usePathname()

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
                        disabled={pathname === "/profile"}
                        hidden={user.role === "admin"}
                        onClick={() => {
                            editBtnHandler()
                        }}
                        className="text-sm md:text-xl w-20 md:w-28"
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

interface Path {
    title: string
    path: string
    visiblility: "all" | "unauthenticated" | "authenticated" | UserRole
}

interface NavItem extends Path {
    id: string
}

// All available routes
const routes: Path[] = [
    { title: "Login", path: "/login", visiblility: "unauthenticated" },
    {
        title: "Parent Portal",
        path: "/parent-portal",
        visiblility: UserRole.GUARDIAN,
    },
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
