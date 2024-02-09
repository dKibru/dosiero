<script>
    import ServerMascot from "./ServerMascot.svg";
    import FileExplorer from "./lib/FileExplorer.svelte";
    import "./lib/TailwindCSS.svelte";
</script>

<section class="body-font text-gray-600">
    <div
        class="container mx-auto flex flex-col items-center justify-center px-5 py-24"
    >
        <div class="inline-block text-center lg:w-2/3">
            <img
                src={ServerMascot}
                class="align-center h-20 text-center"
                alt="upcat mascot"
            />
            <span
                class="title-font mb-4 text-3xl font-medium text-gray-900 sm:text-4xl"
            >
                <h1>Dosiero : updog on steroids</h1>
            </span>
        </div>
        <div class="text-center lg:w-full">
            <div class="mb-10">
                <FileExplorer />
            </div>
        </div>
    </div>
</section>
