<script>
  import Head from "./Head.svelte";
  import Foot from "./Foot.svelte";

  export let blog = {};
</script>

<Head />

<div class="container">
  <h1 class="display-1 mb-4 text-center">{blog.title}</h1>
  <div class="row">
    {@html blog.content}
  </div>
</div>

<Foot />
