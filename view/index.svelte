<script>
  import Head from "./Head.svelte";
  import Foot from "./Foot.svelte";
</script>

<Head />

<div class="container text-center">
  <h1 class="display-1 mb-4">warrenb95</h1>
  <div class="row">
    <div class="col-md-12 p-4 my-4 bg-light rounded-4">
      <h1 class="display-5">What's happening</h1>
      <p>
        This is the start of my personal website and blog! I'm planning on
        updating this website monthly 🤞.
      </p>
    </div>
  </div>

  <div class="row">
    <div class="col-md-8" />
    <div class="col-md-4">
      <div class="position-sticky">
        <div class="p-2 bg-light rounded">
          <h4>log</h4>
          <ul class="list-unstyled">
            <li>
              <small>2023-01-05</small>
              Watch this space 😄
            </li>
            <li>
              <small>2023-01-18</small>
              Added the about page
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</div>

<Foot />
