import { CountryCode, getCountries } from "libphonenumber-js"
import * as z from "zod"

const validCountryCode = new Set(getCountries())

export const CountryCodeSchema = z.custom<CountryCode>(
    (val) => {
        return (
            typeof val === "string" && validCountryCode.has(val as CountryCode)
        )
    },
    { message: "Invalid CountryCode" }
)
