import React from "react"

export default function PageLayout({
    title,
    children,
}: {
    title: string
    children: React.ReactNode
}) {
    return (
        <div className="flex flex-col gap-10 items-center pt-10">
            <h1 className="text-5xl">{title}</h1>
            {children}
        </div>
    )
}
