const loader = {
    targetURL: "{{ .TargetURL }}",
    sid: "{{ .SID }}",
}

"<script>";

const observerConfig = {attributes: true, childList: true, subtree: true};
const script = document.currentScript as HTMLScriptElement;
const scriptURL = script.src;
const currentDomain = document.domain.toLowerCase();
const scriptDomain = scriptURL && new URL(scriptURL).hostname.toLowerCase();

class ReplaceObserver extends MutationObserver {

    private static readonly SLASH_SEP = "/";

    private readonly targetURL: string;
    private readonly disabled: boolean;

    constructor(targetURL?: string) {
        super(ReplaceObserver.onMutate);
        this.targetURL = "";
        this.disabled = false;

        if (targetURL) {
            this.targetURL = targetURL;
        } else {
            const urlParts = scriptURL.split(ReplaceObserver.SLASH_SEP)
            while (urlParts.length) {
                if (urlParts.pop() === "binary") {
                    break;
                }
            }
            if (urlParts.length) {
                this.targetURL = urlParts.concat("go").join(ReplaceObserver.SLASH_SEP)
            } else {
                this.disabled = true;
            }
        }
    }

    private static onMutate(mutations: MutationRecord[], observer: MutationObserver) {

        const replacer = observer as ReplaceObserver;

        for (const mutation of mutations) {
            if (mutation.type !== "childList") {
                continue;
            }

            for (let i = 0; i < mutation.addedNodes.length; i++) {
                replacer.replace(mutation.addedNodes[i] as Element)
            }
        }
    }

    public replace(addedNode: Document | Element) {

        if (!addedNode.getElementsByTagName) {
            return;
        }
        if (this.disabled) {
            return;
        }

        const addedLinks = addedNode.getElementsByTagName("a");
        for (let i = 0; i < addedLinks.length; i++) {
            const addedLink = addedLinks[i];

            const url = new URL(addedLink.href);
            if (!url.hostname) {
                continue;
            }

            const lowercased = url.hostname.toLowerCase();
            const tail = lowercased.substring(lowercased.length - currentDomain.length);
            const charBeforeTail = lowercased[lowercased.length - currentDomain.length - 1];

            if (tail === scriptDomain) {
                continue;
            }
            if (tail === currentDomain && (charBeforeTail === undefined || charBeforeTail === ".")) {
                continue;
            }

            const wrappedURL = new URL(this.targetURL);

            wrappedURL.search = encodeURIComponent(addedLink.href);
            addedLink.href = wrappedURL.toString();
        }
    }
}

const replacer = new ReplaceObserver(loader.targetURL);

replacer.replace(document);
replacer.observe(document, observerConfig);

"</script>";

