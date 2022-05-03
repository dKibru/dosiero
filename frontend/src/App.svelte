<script>
  import logo from './assets/pirate-cat.png'
  import ServerMascot from './ServerMascot.svg';
  import Counter from './lib/Counter.svelte'
  import Filelist from './lib/Filelist.svelte'
  import './lib/TailwindCSS.svelte'

      import axios from "axios"

let base_dir = ""
let directory_items = []



axios.get(import.meta.env.VITE_UPCAT_URL ? `${import.meta.env.VITE_UPCAT_URL}/api/v1/meta` : `/api/v1/meta`)
  .then(function (response) {
    base_dir = response.data.base_directory
    directory_items = response.data.directory_items
    console.log(response.data.directory_items);
  })
  .catch(function (error) {
    console.log(error);
  })
  .then(function () {
  });
  
</script>

<section class="text-gray-600 body-font">
  <div class="container mx-auto flex px-5 py-24 items-center justify-center flex-col">
    
    <div class="text-center lg:w-2/3 inline-block ">
      <img src={ServerMascot} class="h-20 text-center align-center" />
      <span  class="title-font sm:text-4xl text-3xl mb-4 font-medium text-gray-900"> 
        <h1>Dosiero : updog on steroids</h1>
      </span>
      <!-- <p class="mb-8 leading-relaxed">DIR = {base_dir}</p> -->
    </div>
    <div class="text-center lg:w-full">
      <div class="mb-10">
        <Filelist directory_items={directory_items} base_dir={base_dir} />
        <!-- <Counter id="0" /> -->
      </div>
      <!-- <div class="flex justify-center">
        <button class="inline-flex text-white bg-indigo-500 border-0 py-2 px-6 focus:outline-none hover:bg-indigo-600 rounded text-lg">Button</button>
        <button class="ml-4 inline-flex text-gray-700 bg-gray-100 border-0 py-2 px-6 focus:outline-none hover:bg-gray-200 rounded text-lg">Button</button>
      </div> -->
    </div>
  </div>
</section>