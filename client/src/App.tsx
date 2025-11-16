import { Toaster } from "sonner";
import { useQueryState } from "nuqs";
import ShortenTab from "./components/tabs/ShortenTab";
import SuccessTab from "./components/tabs/SuccessTab";
import StatsTab from "./components/tabs/StatsTab";

export default function App() {
    const [tab, setTab] = useQueryState("tab", {
        defaultValue: "shorten",
        history: "push",
        scroll: true,
        parse: (value) => (["shorten", "success", "stats"].includes(value) ? value : "shorten")
    });
    const [shortLinkId, setShortLinkId] = useQueryState("shortLinkId", {
        defaultValue: "",
        history: "push"
    });

    function displayTab() {
        switch (tab) {
            case "shorten":
                return <ShortenTab setShortLinkFn={setShortLinkId} setTabFn={setTab} />;

            case "success":
                return <SuccessTab shortLinkId={shortLinkId} setTabFn={setTab} />;

            case "stats":
                return <StatsTab shortLinkId={shortLinkId} setTabFn={setTab} />;
        }
    }

    return (
        <main className="h-dvh flex flex-col justify-center items-center p-2">
            {displayTab()}
            <Toaster theme="dark" richColors />
        </main>
    );
}
