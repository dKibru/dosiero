<script>
      import axios from "axios"
export let directory_items = [];
export let base_dir = "";
console.log("directory_items", directory_items)

function handleClickDir(event) {
  console.log(event.srcElement.id);
  var item = event.srcElement.id
  axios.get(import.meta.env.VITE_UPCAT_URL +'api/v1/meta?d='+base_dir+"&c="+item)
    .then(function (response) {
      base_dir = response.data.base_directory
      directory_items = response.data.directory_items
    })
    .catch(function (error) {
      console.log(error);
    })
    .then(function () {
    });
}

function handleClickFile(event) {
  console.log(event.srcElement.id);
  var item = event.srcElement.id

    axios({
        url: import.meta.env.VITE_UPCAT_URL +'api/v1/download?d='+base_dir+"&c="+item, 
        method: 'GET',
        responseType: 'blob', 
    })
    .then(function (response) {
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', event.srcElement.id); //or any other extension
      document.body.appendChild(link);
      link.click();
    })
    .catch(function (error) {
      console.log(error);
    })
    .then(function () {
    });
}


</script>
<section class="py-1 bg-blueGray-50">
<div class="w-full xl:w-8/12 mb-12 xl:mb-0 px-4 mx-auto mt-24">
  <div class="relative flex flex-col min-w-0 break-words bg-white w-full mb-6 shadow-lg rounded ">
    <div class="rounded-t mb-0 px-4 py-3 border-0">
      <div class="flex flex-wrap items-center">
        <div class="relative w-full px-4 max-w-full flex-grow flex-1">
          <h3 class="font-semibold text-base text-blueGray-700">DIR = {base_dir}</h3>
        </div>
        <div class="relative w-full px-4 max-w-full flex-grow flex-1 text-right">
          <button class="bg-indigo-500 text-white active:bg-indigo-600 text-xs font-bold uppercase px-3 py-1 rounded outline-none focus:outline-none mr-1 mb-1 ease-linear transition-all duration-150" type="button">See all</button>
        </div>
      </div>
    </div>

    <div class="block w-full overflow-x-auto">
      <table class="items-center bg-transparent w-full border-collapse ">
        <thead>
          <tr>
            <th class="border-t-0 px-6 align-middle border-l-0 border-r-0 text-md whitespace-nowrap p-4 text-left text-black-300">
                          Name
                        </th>
                        <th class="border-t-0 px-6 align-middle border-l-0 border-r-0 text-md whitespace-nowrap p-4 ">
                          Size
                        </th>
                        <th class="border-t-0 px-6 align-middle border-l-0 border-r-0 text-md whitespace-nowrap p-4 ">
                          Last Modified
                        </th>
          </tr>
        </thead>

        <tbody>
          {#each directory_items as ditem}
         


          <tr>
            {#if ditem.IsDir}
            <th on:click={handleClickDir} id={ditem.Id} class="cursor-pointer border-t-0 px-6 align-middle border-l-0 border-r-0 text-xs whitespace-nowrap p-4 text-left text-blue-300 ">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 inline-block fill-current	" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
              </svg>
              {ditem.Name}
            </th>
            {:else}
            <th on:click={handleClickFile} id={ditem.Id} class="cursor-pointer border-t-0 px-6 align-middle border-l-0 border-r-0 text-xs whitespace-nowrap p-4 text-left text-blueGray-300 ">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 inline-block" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>  
              {ditem.Name}
                </th>
              {/if}
            <td class="border-t-0 px-6 align-middle border-l-0 border-r-0 text-xs whitespace-nowrap p-4 ">
              {ditem.Size ? ditem.Size : ""}
            </td>
            <td class="border-t-0 px-6 align-center border-l-0 border-r-0 text-xs whitespace-nowrap p-4">
              {ditem.ModTime}
            </td>
          </tr>
          {/each}
          
        </tbody>

      </table>
    </div>
  </div>
</div>
<footer class="relative pt-8 pb-6 mt-16">
  <div class="container mx-auto px-4">
    <div class="flex flex-wrap items-center md:justify-between justify-center">
      <div class="w-full md:w-6/12 px-4 mx-auto text-center">
        <div class="text-sm text-blueGray-500 font-semibold py-1">
          Made by <a href="https://github.com/dkibru">dkibru</a>
        </div>
      </div>
    </div>
  </div>
</footer>
</section>