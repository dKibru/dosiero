<script>
    import axios from "axios";
    let directory_items = [];
    let base_dir = "";
    const baseApi = import.meta.env.VITE_UPCAT_URL ?? "";

    function getDirContent(dir, file) {
        axios
            .post(`${baseApi}/api/meta`, { dir, file })
            .then(function (response) {
                base_dir = response.data.base_directory;
                directory_items = response.data.directory_items;
            })
            .catch(function (error) {
                console.log(error);
            })
            .then(function () {});
    }
    function handleClickDir(event) {
        var item = event.srcElement.id;
        getDirContent(base_dir, item);
    }

    function handleClickFile(event) {
        var item = event.srcElement.id;
        axios
            .post(
                `${baseApi}/api/download`,
                { dir: base_dir, file: item },
                {
                    responseType: "blob",
                },
            )
            .then(function (response) {
                const url = window.URL.createObjectURL(
                    new Blob([response.data]),
                );
                const link = document.createElement("a");
                link.href = url;
                link.setAttribute("download", event.srcElement.id); //or any other extension
                document.body.appendChild(link);
                link.click();
            })
            .catch(function (error) {
                console.log(error);
            })
            .then(function () {});
    }

    getDirContent("", "");
</script>

<div class="mb-6 w-full max-w-full flex-shrink px-4">
    <div
        class="h-full rounded-lg bg-white px-10 py-6 shadow-lg dark:bg-gray-800"
    >
        <div class="relative">
            <div class="-mx-4 flex flex-row flex-wrap">
                <div class="w-full max-w-full flex-shrink px-4">
                    <div class="relative">
                        <h3
                            class="mb-3 text-lg font-semibold text-gray-600 dark:text-gray-200"
                        >
                            {base_dir}
                        </h3>
                        <!-- content filter start -->
                        <table
                            class="w-full table-auto text-sm ltr:text-left rtl:text-right"
                        >
                            <thead
                                class="bg-gray-100 dark:bg-gray-900 dark:bg-opacity-40"
                            >
                                <tr>
                                    <th class="px-4 py-3 font-normal"
                                        >File name</th
                                    >
                                    <th
                                        class="hidden px-4 py-3 font-normal md:table-cell"
                                        >Last modified</th
                                    >
                                    <th class="px-4 py-3 font-normal"
                                        >File size</th
                                    >
                                    <th
                                        class="px-4 py-3 text-center font-normal"
                                        >Action</th
                                    >
                                </tr>
                            </thead>
                            <tbody>
                                {#each directory_items as ditem}
                                    <tr>
                                        <td
                                            class="px-4 py-3 text-right font-medium"
                                        >
                                            <div
                                                class="cursor-pointer"
                                                on:click={ditem.IsDir
                                                    ? handleClickDir
                                                    : handleClickFile}
                                                id={ditem.Id}
                                            >
                                                <span
                                                    class="mr-1 text-indigo-500"
                                                >
                                                    {#if ditem.IsDir}
                                                        <svg
                                                            xmlns="http://www.w3.org/2000/svg"
                                                            width="16"
                                                            height="16"
                                                            fill="currentColor"
                                                            class="bi bi-folder-fill inline-block"
                                                            viewBox="0 0 16 16"
                                                        >
                                                            <path
                                                                d="M9.828 3h3.982a2 2 0 0 1 1.992 2.181l-.637 7A2 2 0 0 1 13.174 14H2.825a2 2 0 0 1-1.991-1.819l-.637-7a1.99 1.99 0 0 1 .342-1.31L.5 3a2 2 0 0 1 2-2h3.672a2 2 0 0 1 1.414.586l.828.828A2 2 0 0 0 9.828 3zm-8.322.12C1.72 3.042 1.95 3 2.19 3h5.396l-.707-.707A1 1 0 0 0 6.172 2H2.5a1 1 0 0 0-1 .981l.006.139z"
                                                            ></path>
                                                        </svg>
                                                    {:else}
                                                        <svg
                                                            xmlns="http://www.w3.org/2000/svg"
                                                            width="16"
                                                            height="16"
                                                            fill="currentColor"
                                                            class="bi bi-file-zip inline-block"
                                                            viewBox="0 0 16 16"
                                                        >
                                                            <path
                                                                d="M6.5 7.5a1 1 0 0 1 1-1h1a1 1 0 0 1 1 1v.938l.4 1.599a1 1 0 0 1-.416 1.074l-.93.62a1 1 0 0 1-1.109 0l-.93-.62a1 1 0 0 1-.415-1.074l.4-1.599V7.5zm2 0h-1v.938a1 1 0 0 1-.03.243l-.4 1.598.93.62.93-.62-.4-1.598a1 1 0 0 1-.03-.243V7.5z"
                                                            ></path>
                                                            <path
                                                                d="M2 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V2zm5.5-1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H9v1H8v1h1v1H8v1h1v1H7.5V5h-1V4h1V3h-1V2h1V1z"
                                                            ></path>
                                                        </svg>
                                                    {/if}
                                                </span>{ditem.Name}
                                            </div>
                                        </td>
                                        <td
                                            class="hidden px-4 py-3 font-medium md:table-cell"
                                            >{ditem.ModTime}</td
                                        >
                                        <td class="px-4 py-3 font-medium"
                                            >{ditem.Size ? ditem.Size : ""}</td
                                        >
                                        <td
                                            class="px-4 py-3 text-center font-medium"
                                        >
                                            <a
                                                href="javascript:;"
                                                class="inline-block hover:text-red-500 ltr:mr-2 rtl:ml-2"
                                                title="Delete"
                                            >
                                                <svg
                                                    xmlns="http://www.w3.org/2000/svg"
                                                    width="16"
                                                    height="16"
                                                    fill="currentColor"
                                                    class="bi bi-trash"
                                                    viewBox="0 0 16 16"
                                                >
                                                    <path
                                                        d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"
                                                    ></path>
                                                    <path
                                                        fill-rule="evenodd"
                                                        d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"
                                                    ></path>
                                                </svg>
                                            </a>
                                            <a
                                                href="javascript:;"
                                                class="inline-block hover:text-green-500 ltr:mr-2 rtl:ml-2"
                                                title="Rename"
                                            >
                                                <svg
                                                    xmlns="http://www.w3.org/2000/svg"
                                                    width="16"
                                                    height="16"
                                                    fill="currentColor"
                                                    class="bi bi-pencil-square"
                                                    viewBox="0 0 16 16"
                                                >
                                                    <path
                                                        d="M15.502 1.94a.5.5 0 0 1 0 .706L14.459 3.69l-2-2L13.502.646a.5.5 0 0 1 .707 0l1.293 1.293zm-1.75 2.456-2-2L4.939 9.21a.5.5 0 0 0-.121.196l-.805 2.414a.25.25 0 0 0 .316.316l2.414-.805a.5.5 0 0 0 .196-.12l6.813-6.814z"
                                                    ></path>
                                                    <path
                                                        fill-rule="evenodd"
                                                        d="M1 13.5A1.5 1.5 0 0 0 2.5 15h11a1.5 1.5 0 0 0 1.5-1.5v-6a.5.5 0 0 0-1 0v6a.5.5 0 0 1-.5.5h-11a.5.5 0 0 1-.5-.5v-11a.5.5 0 0 1 .5-.5H9a.5.5 0 0 0 0-1H2.5A1.5 1.5 0 0 0 1 2.5v11z"
                                                    ></path>
                                                </svg>
                                            </a>
                                        </td>
                                    </tr>
                                {/each}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
