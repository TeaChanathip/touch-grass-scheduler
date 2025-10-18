import CountryList from "country-list-with-dial-code-and-flag"
import { getCountries } from "libphonenumber-js"
import { UserGender, UserRole } from "../interfaces/User.interface"
import snakeToTitleCase from "../utils/snakeToTitleCase"

export const genderOptions = Object.keys(UserGender).map((key) => {
    return {
        value: UserGender[key],
        label: snakeToTitleCase(UserGender[key]),
        id: `select-gender-${UserGender[key]}`,
    }
})

export const roleOptions = Object.keys(UserRole)
    .filter((key) => UserRole[key] !== "admin")
    .map((key) => {
        return {
            value: UserRole[key],
            label: snakeToTitleCase(UserRole[key]),
            id: `select-role-${UserRole[key]}`,
        }
    })

export const countryCodesOptions = getCountries()
    .map((countryCode) => {
        const country = CountryList.findOneByCountryCode(countryCode)

        let label = ""
        if (country) {
            label = `(${country.dial_code}) ${country.name}`
        }

        return {
            value: countryCode,
            label: label,
            id: `select-country-code-${countryCode}`,
        }
    })
    .sort((option1, option2) => {
        return option1.value.localeCompare(option2.value)
    })
