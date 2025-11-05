import EmailRoundedIcon from "@mui/icons-material/EmailRounded"
import DraftsRoundedIcon from "@mui/icons-material/DraftsRounded"

import { Dispatch, memo, SetStateAction, useRef, useState } from "react"
import { Badge } from "@mui/material"
import useClickOutside from "../hooks/useClickOutside"

const NotiButton = memo(function NotiButton() {
    // Hooks
    const [isOpen, setOpen] = useState(false)
    const notiBoxRef = useRef<HTMLDivElement | null>(null)
    useClickOutside(notiBoxRef, () => setOpen(false))

    const onClickHandler = () => {
        setOpen((isOpen) => !isOpen)
    }

    return (
        <div className="relative" ref={notiBoxRef}>
            <button
                onClick={() => onClickHandler()}
                className="text-white cursor-pointer"
            >
                {isOpen ? (
                    <DraftsRoundedIcon sx={{ fontSize: 32 }} />
                ) : (
                    <Badge
                        color="error"
                        overlap="circular"
                        badgeContent="10"
                        invisible={false}
                    >
                        <EmailRoundedIcon sx={{ fontSize: 32 }} />
                    </Badge>
                )}
            </button>
            <NotiBox isOpen={isOpen} setOpen={setOpen} />
        </div>
    )
})
export default NotiButton

function NotiBox({
    isOpen,
    setOpen,
}: {
    isOpen: boolean
    setOpen: Dispatch<SetStateAction<boolean>>
}) {
    return (
        <div
            className="fixed left-1/2 -translate-x-1/2 w-[90vw] h-[200px]
                md:absolute md:left-auto md:right-0 md:translate-none md: md:w-[400px]
                bg-prim-green-50 rounded-2xl z-0 mt-3 
                transition-all transition-discrete duration-150 ease-in-out"
            style={
                isOpen
                    ? { opacity: "90%" }
                    : { opacity: "0%", pointerEvents: "none" }
            }
        ></div>
    )
}

function NotiItem() {
    return <></>
}
