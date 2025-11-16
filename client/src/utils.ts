import { AxiosError } from "axios";

export function formatError(error: unknown) {
    return error instanceof AxiosError && error.response?.status === 400
        ? `, link is invalid`
        : error instanceof AxiosError && error.response?.status === 429
        ? ", too many requests. Wait for a few minutes"
        : error instanceof AxiosError
        ? `, ${error.response?.data?.error || "an error occurred"}`
        : "";
}

export function linkShortIdToUrl(shortLinkId: string) {
    return new URL(`/${shortLinkId}`, location.href).toString();
}
