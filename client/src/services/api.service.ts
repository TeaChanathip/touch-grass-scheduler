export class ApiService {
    private timeoutMs = 8000

    private async safeParseJson(res: Response) {
        const text = await res.text()
        if (!text) return null
        try {
            return JSON.parse(text)
        } catch {
            // Not JSON
            return text
        }
    }

    private async request<T, R>(
        method: string,
        url: string,
        options?: {
            params?: Record<string, any>
            body?: T
            timeoutMs?: number
        }
    ): Promise<R> {
        const controller = new AbortController()
        const timeoutMs = options?.timeoutMs ?? this.timeoutMs
        const id = setTimeout(() => controller.abort(), timeoutMs)

        // Build URL with query string if params provided
        let fetchUrl = url
        if (options?.params) {
            const searchParams = new URLSearchParams()
            for (const [k, v] of Object.entries(options.params)) {
                if (v === undefined || v === null) continue
                if (Array.isArray(v)) {
                    v.forEach((item) => searchParams.append(k, String(item)))
                } else {
                    searchParams.set(k, String(v))
                }
            }
            const qs = searchParams.toString()
            if (qs) {
                // append correctly whether url already has query
                fetchUrl += (url.includes("?") ? "&" : "?") + qs
            }
        }

        const headers = new Headers({ "Content-Type": "application/json" })

        try {
            const response = await fetch(fetchUrl, {
                method,
                headers,
                credentials: "include",
                mode: "cors", // Explicitly set CORS mode
                cache: "no-cache", // Prevent caching
                body:
                    options?.body !== undefined && method !== "GET"
                        ? JSON.stringify(options.body)
                        : undefined,
                signal: controller.signal,
            })

            clearTimeout(id)

            const parsed = await this.safeParseJson(response)

            if (!response.ok) {
                // Prefer server-provided message if available
                let errMsg = `HTTP ${response.status}`
                if (
                    parsed &&
                    typeof parsed === "object" &&
                    (parsed as any).error
                ) {
                    errMsg = (parsed as any).error
                } else if (typeof parsed === "string" && parsed.length > 0) {
                    errMsg = parsed
                }

                console.log(errMsg)

                throw new ApiError(errMsg, response.status)
            }

            return parsed as R
        } catch (err) {
            if ((err as any)?.name === "AbortError") {
                throw new Error("Request timed out")
            }
            throw err
        } finally {
            clearTimeout(id)
        }
    }

    async get<R>(
        url: string,
        params?: Record<string, any>,
        timeoutMs?: number
    ) {
        return this.request<undefined, R>("GET", url, { params, timeoutMs })
    }

    async post<T, R>(
        url: string,
        body?: T,
        params?: Record<string, any>,
        timeoutMs?: number
    ) {
        return this.request<T, R>("POST", url, { body, params, timeoutMs })
    }

    async put<T, R>(
        url: string,
        body?: T,
        params?: Record<string, any>,
        timeoutMs?: number
    ) {
        return this.request<T, R>("PUT", url, { body, params, timeoutMs })
    }

    async delete<T, R>(
        url: string,
        body?: T,
        params?: Record<string, any>,
        timeoutMs?: number
    ) {
        return this.request<T, R>("DELETE", url, { body, params, timeoutMs })
    }
}

export class ApiError extends Error {
    status?: number

    constructor(message: string, status?: number) {
        super(message)
        this.name = "ApiError"
        this.status = status
    }
}

export function isApiError(e: unknown): e is ApiError {
    return typeof e === "object" && e !== null && (e as any).name === "ApiError"
}
