import { toast } from "sonner";
import type { Tabs } from "../../types";
import { linkShortIdToUrl } from "../../utils";
import { QRCodeCanvas } from "qrcode.react";

type Props = {
    shortLinkId: string;
    setTabFn: (tabName: Tabs) => void;
};

export default function SuccessTab({ shortLinkId, setTabFn }: Readonly<Props>) {
    const shortUrl = linkShortIdToUrl(shortLinkId);

    async function copyLinkToClipboard() {
        if (!navigator.clipboard) {
            toast.error("Clipboard not available");
            return;
        }

        try {
            await navigator.clipboard.writeText(shortUrl);
            toast.success("Link copied");
        } catch (error) {
            console.error("Failed to copy to clipboard:", error);
            toast.error("Error copying link");
        }
    }

    return (
        <div className="bg-muted border border-border rounded-sm flex flex-col gap-1 p-2">
            <h2 className="text-2xl font-semibold text-center">Your shortened URL</h2>
            <div className="flex gap-1">
                <p className="px-1.5 py-1 border border-border rounded-sm select-all">{shortUrl}</p>
                <button
                    className="px-1.5 py-1 bg-accent rounded-sm shadow-sm transition-all hover:brightness-110 cursor-pointer"
                    onClick={copyLinkToClipboard}
                >
                    Copy
                </button>
            </div>
            <div className="w-full flex justify-center">
                <QRCodeCanvas value={shortUrl} marginSize={1} title={`QR code for ${shortUrl}`} />
            </div>
            <button
                className="px-1.5 py-1 bg-accent rounded-sm shadow-sm transition-all hover:brightness-110 cursor-pointer"
                onClick={() => setTabFn("stats")}
            >
                View link stats
            </button>
            <button
                className="px-1.5 py-1 bg-accent rounded-sm shadow-sm transition-all hover:brightness-110 cursor-pointer"
                onClick={() => setTabFn("shorten")}
            >
                Shorten another URL
            </button>
        </div>
    );
}
