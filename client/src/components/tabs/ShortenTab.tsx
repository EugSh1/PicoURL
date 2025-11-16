import { useState, type ChangeEvent, type KeyboardEvent } from "react";
import { toast } from "sonner";
import type { ShortenLinkResponse, Tabs } from "../../types";
import { formatError } from "../../utils";
import axios from "axios";
import { urlRegexp } from "../../constants";

type Props = {
    setShortLinkFn: (url: string) => void;
    setTabFn: (tabName: Tabs) => void;
};

export default function ShortenTab({ setShortLinkFn: setShortLinkId, setTabFn }: Readonly<Props>) {
    const [linkToShorten, setLinkToShorten] = useState<string>("");
    const [isLoading, setIsLoading] = useState<boolean>(false);

    async function shortenLink() {
        if (isLoading) return;

        const trimmedUrl = linkToShorten.trim();

        if (!urlRegexp.test(trimmedUrl)) {
            toast.error("Link is invalid");
            return;
        }

        try {
            setIsLoading(true);

            const response = await axios.post<ShortenLinkResponse>("/api/shorten", {
                url: trimmedUrl
            });
            setShortLinkId(response.data.shortUrlId);
            setTabFn("success");
            toast.success("Link shortened");
            setLinkToShorten("");
        } catch (error) {
            console.error("Error shortening link", error);
            toast.error(`Error shortening link${formatError(error)}`);
        } finally {
            setIsLoading(false);
        }
    }

    function handleKeyDown(event: KeyboardEvent<HTMLInputElement>) {
        if (event.key === "Enter") {
            shortenLink();
        }
    }

    return (
        <>
            <h1 className="text-4xl font-semibold mb-1">PicoUrl</h1>
            <div className="flex gap-1">
                <input
                    aria-label="URL to shorten"
                    type="url"
                    placeholder="Enter the link here"
                    className="px-1.5 py-1 bg-muted border border-border rounded-sm focus-visible:outline user-valid:outline-accent"
                    value={linkToShorten}
                    onChange={(event: ChangeEvent<HTMLInputElement>) =>
                        setLinkToShorten(event.target.value)
                    }
                    onKeyDown={handleKeyDown}
                />
                <button
                    className="px-1.5 py-1 bg-accent rounded-sm shadow-sm transition-all hover:brightness-110 cursor-pointer"
                    onClick={shortenLink}
                    disabled={isLoading}
                >
                    {isLoading ? "Shortening..." : "Shorten"}
                </button>
            </div>
        </>
    );
}
