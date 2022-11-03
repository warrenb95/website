<script>
  import Head from './Head.svelte'
  import Foot from './Foot.svelte'
  export let blogs = []
</script>

<Head />

<div class="containter px-5">

    <h1>Blog Index</h1>

    <table class="table">
      {#if blogs.length > 0}
        <thead>
          {#each Object.keys(blogs[0]) as key}
            <th scope="col">{key}</th>
          {/each}
        </thead>
      {/if}
            <tbody>
      {#each blogs as blog}
        <tr>
          {#each Object.keys(blog) as key}
            {#if key.toLowerCase() === "id"}
              <td><a href={`/${blog.id || 0}`}>{blog[key]}</a></td>
            {:else}
              <td>{blog[key]}</td>
            {/if}
          {/each}
        </tr>
      {/each}
            </tbody>
    </table>

    <br />

    <a href={`/new`}>New Blog</a>

</div>

<Foot />
