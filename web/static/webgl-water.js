/**
 * WebGL Water Tutorial - JavaScript Frontend
 * Go Port of the original Rust/WASM implementation
 */

class WebGLWaterApp {
  constructor() {
    this.canvas = null;
    this.gl = null;
    this.shaders = {};
    this.programs = {};
    this.meshes = {};
    this.textures = {};
    this.framebuffers = {};

    // Application state
    this.state = {
      clock: 0,
      camera: {
        position: [0, 5, 10],
        viewMatrix: null,
      },
      water: {
        reflectivity: 0.6,
        fresnelStrength: 2.0,
        waveSpeed: 0.03,
        useReflection: true,
        useRefraction: true,
      },
      scenery: true,
    };

    // Constants
    this.CANVAS_WIDTH = 1200;
    this.CANVAS_HEIGHT = 800;
    this.WATER_TILE_Y_POS = 0.0;
    this.REFLECTION_TEXTURE_WIDTH = 320;
    this.REFLECTION_TEXTURE_HEIGHT = 180;
    this.REFRACTION_TEXTURE_WIDTH = 1280;
    this.REFRACTION_TEXTURE_HEIGHT = 720;

    // WebSocket connection
    this.ws = null;

    // Mouse interaction state
    this.mousePressed = false;
    this.lastMouseX = 0;
    this.lastMouseY = 0;

    // Initialize
    this.init();
  }

  async init() {
    try {
      console.log("🌊 Starting WebGL Water initialization...");
      this.setupCanvas();
      console.log("✅ Canvas setup complete");
      this.setupWebGL();
      console.log("✅ WebGL context created");
      await this.loadShaders();
      console.log("✅ Shaders loaded");
      await this.loadAssets();
      console.log("✅ Assets loaded");
      this.setupFramebuffers();
      console.log("✅ Framebuffers setup");
      this.setupEventHandlers();
      console.log("✅ Event handlers setup");
      this.connectWebSocket();
      console.log("✅ WebSocket connected");

      // Start render loop
      this.lastTime = Date.now();
      this.render();

      console.log("🎉 WebGL Water application initialized successfully");
    } catch (error) {
      console.error("❌ Failed to initialize WebGL Water application:", error);
      console.error("Stack trace:", error.stack);
    }
  }

  setupCanvas() {
    this.canvas = document.getElementById("canvas");
    if (!this.canvas) {
      throw new Error("Canvas element not found");
    }

    // Set canvas size
    this.canvas.width = this.CANVAS_WIDTH;
    this.canvas.height = this.CANVAS_HEIGHT;
  }

  setupWebGL() {
    console.log("🔧 Setting up WebGL context...");
    this.gl =
      this.canvas.getContext("webgl") ||
      this.canvas.getContext("experimental-webgl");
    if (!this.gl) {
      console.error("❌ WebGL not supported by this browser");
      throw new Error("WebGL not supported");
    }
    console.log("✅ WebGL context created successfully");

    const gl = this.gl;

    // Enable depth testing
    gl.enable(gl.DEPTH_TEST);
    gl.depthFunc(gl.LEQUAL);

    // Enable blending for transparency
    gl.enable(gl.BLEND);
    gl.blendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);

    // Set clear color (sky blue)
    gl.clearColor(0.53, 0.8, 0.98, 1.0);

    // Set viewport
    gl.viewport(0, 0, this.CANVAS_WIDTH, this.CANVAS_HEIGHT);

    // Get extensions
    this.vaoExt = gl.getExtension("OES_vertex_array_object");
    this.depthTextureExt = gl.getExtension("WEBGL_depth_texture");

    if (!this.vaoExt) {
      console.warn("⚠️ VAO extension not available");
    } else {
      console.log("✅ VAO extension loaded");
    }

    if (!this.depthTextureExt) {
      console.warn("⚠️ Depth texture extension not available");
    } else {
      console.log("✅ Depth texture extension loaded");
    }
  }

  async loadShaders() {
    const shaderNames = [
      "water-vertex",
      "water-fragment",
      "mesh-vertex",
      "mesh-fragment",
      "textured-quad-vertex",
      "textured-quad-fragment",
    ];

    // Load shader sources
    for (const name of shaderNames) {
      try {
        const response = await fetch(`/shaders/${name}.glsl`);
        if (!response.ok) {
          throw new Error(`Failed to load shader: ${name}`);
        }
        this.shaders[name] = await response.text();
        console.log(`✅ Loaded shader: ${name}`);
      } catch (error) {
        console.error(`❌ Error loading shader ${name}:`, error);
        throw error;
      }
    }

    // Compile and link shader programs
    this.programs.water = this.createProgram("water-vertex", "water-fragment");
    this.programs.mesh = this.createProgram("mesh-vertex", "mesh-fragment");
    this.programs.quad = this.createProgram(
      "textured-quad-vertex",
      "textured-quad-fragment",
    );
  }

  createProgram(vertexShaderName, fragmentShaderName) {
    const gl = this.gl;

    const vertexShader = this.compileShader(
      gl.VERTEX_SHADER,
      this.shaders[vertexShaderName],
    );
    const fragmentShader = this.compileShader(
      gl.FRAGMENT_SHADER,
      this.shaders[fragmentShaderName],
    );

    const program = gl.createProgram();
    gl.attachShader(program, vertexShader);
    gl.attachShader(program, fragmentShader);
    gl.linkProgram(program);

    if (!gl.getProgramParameter(program, gl.LINK_STATUS)) {
      const error = gl.getProgramInfoLog(program);
      gl.deleteProgram(program);
      throw new Error(`Program linking failed: ${error}`);
    }

    // Get attribute and uniform locations
    program.attribLocations = {};
    program.uniformLocations = {};

    const numAttribs = gl.getProgramParameter(program, gl.ACTIVE_ATTRIBUTES);
    for (let i = 0; i < numAttribs; i++) {
      const attrib = gl.getActiveAttrib(program, i);
      program.attribLocations[attrib.name] = gl.getAttribLocation(
        program,
        attrib.name,
      );
    }

    const numUniforms = gl.getProgramParameter(program, gl.ACTIVE_UNIFORMS);
    for (let i = 0; i < numUniforms; i++) {
      const uniform = gl.getActiveUniform(program, i);
      program.uniformLocations[uniform.name] = gl.getUniformLocation(
        program,
        uniform.name,
      );
    }

    return program;
  }

  compileShader(type, source) {
    const gl = this.gl;
    const shader = gl.createShader(type);

    gl.shaderSource(shader, source);
    gl.compileShader(shader);

    if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
      const error = gl.getShaderInfoLog(shader);
      gl.deleteShader(shader);
      throw new Error(`Shader compilation failed: ${error}`);
    }

    return shader;
  }

  async loadAssets() {
    // Load meshes
    const meshResponse = await fetch("/api/meshes");
    const meshData = await meshResponse.json();

    for (const meshName of meshData.meshes) {
      console.log(`📦 Loading mesh: ${meshName}`);
      const response = await fetch(`/api/meshes/${meshName}`);
      if (!response.ok) {
        throw new Error(`Failed to load mesh ${meshName}: ${response.status}`);
      }
      const mesh = await response.json();
      this.meshes[meshName] = this.createMeshBuffers(mesh);
      console.log(`✅ Mesh loaded: ${meshName} (${mesh.vertexCount} vertices)`);
    }

    // Load textures
    const textureNames = [
      { name: "dudvmap", file: "dudvmap.png" },
      { name: "normalmap", file: "normalmap.png" },
      { name: "stone", file: "stone-texture.png" },
    ];
    for (const texture of textureNames) {
      console.log(`🖼️ Loading texture: ${texture.name}`);
      await this.loadTexture(texture.name, `/assets/${texture.file}`);
    }
  }

  createMeshBuffers(meshData) {
    const gl = this.gl;

    const mesh = {
      vertexBuffer: gl.createBuffer(),
      normalBuffer: gl.createBuffer(),
      texCoordBuffer: gl.createBuffer(),
      indexBuffer: gl.createBuffer(),
      vertexCount: meshData.vertexCount,
      indexCount: meshData.indices.length,
    };

    // Vertex positions
    gl.bindBuffer(gl.ARRAY_BUFFER, mesh.vertexBuffer);
    gl.bufferData(
      gl.ARRAY_BUFFER,
      new Float32Array(meshData.vertices),
      gl.STATIC_DRAW,
    );

    // Normals
    gl.bindBuffer(gl.ARRAY_BUFFER, mesh.normalBuffer);
    gl.bufferData(
      gl.ARRAY_BUFFER,
      new Float32Array(meshData.normals),
      gl.STATIC_DRAW,
    );

    // Texture coordinates
    gl.bindBuffer(gl.ARRAY_BUFFER, mesh.texCoordBuffer);
    gl.bufferData(
      gl.ARRAY_BUFFER,
      new Float32Array(meshData.texCoords),
      gl.STATIC_DRAW,
    );

    // Indices
    gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.indexBuffer);
    gl.bufferData(
      gl.ELEMENT_ARRAY_BUFFER,
      new Uint16Array(meshData.indices),
      gl.STATIC_DRAW,
    );

    return mesh;
  }

  async loadTexture(name, url) {
    return new Promise((resolve, reject) => {
      const gl = this.gl;
      const texture = gl.createTexture();
      const image = new Image();

      image.onload = () => {
        gl.bindTexture(gl.TEXTURE_2D, texture);
        gl.texImage2D(
          gl.TEXTURE_2D,
          0,
          gl.RGBA,
          gl.RGBA,
          gl.UNSIGNED_BYTE,
          image,
        );

        // Generate mipmaps for better filtering
        gl.generateMipmap(gl.TEXTURE_2D);

        // Set texture parameters
        gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR);
        gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
        gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT);
        gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT);

        this.textures[name] = texture;
        console.log(
          `✅ Texture loaded: ${name} (${image.width}x${image.height})`,
        );
        resolve(texture);
      };

      image.onerror = () => {
        console.error(`❌ Failed to load texture: ${url}`);
        reject(new Error(`Failed to load texture: ${url}`));
      };

      image.src = url;
    });
  }

  setupFramebuffers() {
    this.framebuffers.reflection = this.createFramebuffer(
      this.REFLECTION_TEXTURE_WIDTH,
      this.REFLECTION_TEXTURE_HEIGHT,
    );

    this.framebuffers.refraction = this.createFramebuffer(
      this.REFRACTION_TEXTURE_WIDTH,
      this.REFRACTION_TEXTURE_HEIGHT,
    );
  }

  createFramebuffer(width, height) {
    const gl = this.gl;

    const framebuffer = gl.createFramebuffer();
    gl.bindFramebuffer(gl.FRAMEBUFFER, framebuffer);

    // Color texture
    const colorTexture = gl.createTexture();
    gl.bindTexture(gl.TEXTURE_2D, colorTexture);
    gl.texImage2D(
      gl.TEXTURE_2D,
      0,
      gl.RGB,
      width,
      height,
      0,
      gl.RGB,
      gl.UNSIGNED_BYTE,
      null,
    );
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
    gl.framebufferTexture2D(
      gl.FRAMEBUFFER,
      gl.COLOR_ATTACHMENT0,
      gl.TEXTURE_2D,
      colorTexture,
      0,
    );

    // Depth texture
    const depthTexture = gl.createTexture();
    gl.bindTexture(gl.TEXTURE_2D, depthTexture);
    gl.texImage2D(
      gl.TEXTURE_2D,
      0,
      gl.DEPTH_COMPONENT,
      width,
      height,
      0,
      gl.DEPTH_COMPONENT,
      gl.UNSIGNED_SHORT,
      null,
    );
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
    gl.framebufferTexture2D(
      gl.FRAMEBUFFER,
      gl.DEPTH_ATTACHMENT,
      gl.TEXTURE_2D,
      depthTexture,
      0,
    );

    if (gl.checkFramebufferStatus(gl.FRAMEBUFFER) !== gl.FRAMEBUFFER_COMPLETE) {
      throw new Error("Framebuffer is not complete");
    }

    gl.bindFramebuffer(gl.FRAMEBUFFER, null);

    return {
      framebuffer,
      colorTexture,
      depthTexture,
      width,
      height,
    };
  }

  setupEventHandlers() {
    // Mouse controls
    this.canvas.addEventListener("mousedown", this.onMouseDown.bind(this));
    this.canvas.addEventListener("mouseup", this.onMouseUp.bind(this));
    this.canvas.addEventListener("mousemove", this.onMouseMove.bind(this));
    this.canvas.addEventListener("wheel", this.onWheel.bind(this));

    // UI controls
    this.setupUIControls();
  }

  setupUIControls() {
    const controls = {
      reflectivity: (value) =>
        this.updateWaterProperty("reflectivity", parseFloat(value)),
      fresnel: (value) =>
        this.updateWaterProperty("fresnelStrength", parseFloat(value)),
      "wave-speed": (value) =>
        this.updateWaterProperty("waveSpeed", parseFloat(value)),
      "use-reflection": (value) =>
        this.updateWaterProperty("useReflection", value),
      "use-refraction": (value) =>
        this.updateWaterProperty("useRefraction", value),
      "show-scenery": (value) => this.updateScenery(value),
    };

    for (const [id, handler] of Object.entries(controls)) {
      const element = document.getElementById(id);
      if (element) {
        if (element.type === "range") {
          element.addEventListener("input", (e) => {
            handler(e.target.value);
            const valueElement = document.getElementById(`${id}-value`);
            if (valueElement) {
              valueElement.textContent = e.target.value;
            }
          });
        } else if (element.type === "checkbox") {
          element.addEventListener("change", (e) => {
            handler(e.target.checked);
          });
        }
      }
    }
  }

  onMouseDown(event) {
    this.mousePressed = true;
    this.lastMouseX = event.clientX;
    this.lastMouseY = event.clientY;

    this.sendCameraUpdate({
      mouseDown: { x: event.clientX, y: event.clientY },
    });
  }

  onMouseUp(event) {
    this.mousePressed = false;

    this.sendCameraUpdate({
      mouseUp: true,
    });
  }

  onMouseMove(event) {
    if (this.mousePressed) {
      this.sendCameraUpdate({
        mouseMove: { x: event.clientX, y: event.clientY },
      });
    }

    this.lastMouseX = event.clientX;
    this.lastMouseY = event.clientY;
  }

  onWheel(event) {
    event.preventDefault();
    const delta = event.deltaY * 0.01;

    this.sendCameraUpdate({
      zoom: delta,
    });
  }

  async updateWaterProperty(property, value) {
    const update = {};
    update[property] = value;

    try {
      await fetch("/api/state/water", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(update),
      });
    } catch (error) {
      console.error("Failed to update water property:", error);
    }
  }

  async updateScenery(show) {
    // For now, just update local state
    // In a full implementation, this would be sent to the server
    this.state.scenery = show;
  }

  async sendCameraUpdate(update) {
    try {
      await fetch("/api/state/camera", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(update),
      });
    } catch (error) {
      console.error("Failed to update camera:", error);
    }
  }

  connectWebSocket() {
    const protocol = location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${location.host}/ws`;

    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      console.log("WebSocket connected");
    };

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "state_update") {
          this.state = { ...this.state, ...data };
        }
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
      }
    };

    this.ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    this.ws.onclose = () => {
      console.log("WebSocket disconnected, attempting to reconnect...");
      setTimeout(() => {
        this.connectWebSocket();
      }, 2000);
    };
  }

  render() {
    const currentTime = Date.now();
    const deltaTime = currentTime - this.lastTime;
    this.lastTime = currentTime;

    // Update local clock for animation
    this.state.clock = currentTime;

    const gl = this.gl;

    if (!gl) {
      console.error("❌ WebGL context lost during render");
      return;
    }

    // Clear the main framebuffer
    gl.bindFramebuffer(gl.FRAMEBUFFER, null);
    gl.viewport(0, 0, this.CANVAS_WIDTH, this.CANVAS_HEIGHT);
    gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);

    // Render refraction to framebuffer
    this.renderRefraction();

    // Render reflection to framebuffer
    this.renderReflection();

    // Render main scene
    this.renderMainScene();

    // Continue render loop
    requestAnimationFrame(this.render.bind(this));
  }

  renderRefraction() {
    if (!this.state.water.useRefraction) return;

    const gl = this.gl;
    const fb = this.framebuffers.refraction;

    gl.bindFramebuffer(gl.FRAMEBUFFER, fb.framebuffer);
    gl.viewport(0, 0, fb.width, fb.height);
    gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);

    // Render scene below water (with clipping plane)
    const clipPlane = [0, -1, 0, this.WATER_TILE_Y_POS];
    this.renderMeshes(clipPlane, false);
  }

  renderReflection() {
    if (!this.state.water.useReflection) return;

    const gl = this.gl;
    const fb = this.framebuffers.reflection;

    gl.bindFramebuffer(gl.FRAMEBUFFER, fb.framebuffer);
    gl.viewport(0, 0, fb.width, fb.height);
    gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);

    // Render scene above water (with clipping plane, mirrored)
    const clipPlane = [0, 1, 0, -this.WATER_TILE_Y_POS];
    this.renderMeshes(clipPlane, true);
  }

  renderMainScene() {
    const gl = this.gl;

    gl.bindFramebuffer(gl.FRAMEBUFFER, null);
    gl.viewport(0, 0, this.CANVAS_WIDTH, this.CANVAS_HEIGHT);

    // Render water
    this.renderWater();

    // Render scene meshes
    const clipPlane = [0, 1, 0, 1000000]; // No clipping
    this.renderMeshes(clipPlane, false);

    // Render debug views (small previews of the framebuffer textures)
    this.renderDebugViews();
  }

  renderWater() {
    const gl = this.gl;
    const program = this.programs.water;
    const mesh = this.meshes.water_plane;

    if (!program || !mesh) return;

    gl.useProgram(program);

    // Bind vertex attributes
    this.bindMeshAttributes(program, mesh);

    // Set uniforms
    const perspectiveMatrix = this.getPerspectiveMatrix();
    const viewMatrix = this.getViewMatrix();
    const modelMatrix = this.getIdentityMatrix();
    const cameraPos = this.state.camera.position;

    gl.uniformMatrix4fv(
      program.uniformLocations.perspective,
      false,
      perspectiveMatrix,
    );
    gl.uniformMatrix4fv(program.uniformLocations.view, false, viewMatrix);
    gl.uniformMatrix4fv(program.uniformLocations.model, false, modelMatrix);
    gl.uniform3fv(program.uniformLocations.cameraPos, cameraPos);

    // Water-specific uniforms
    const dudvOffset = (this.state.clock / 1000.0) * this.state.water.waveSpeed;
    gl.uniform1f(program.uniformLocations.dudvOffset, dudvOffset % 1.0);
    gl.uniform1f(
      program.uniformLocations.waterReflectivity,
      this.state.water.reflectivity,
    );
    gl.uniform1f(
      program.uniformLocations.fresnelStrength,
      this.state.water.fresnelStrength,
    );

    // Bind textures
    this.bindTexture(gl.TEXTURE0, this.framebuffers.refraction.colorTexture);
    gl.uniform1i(program.uniformLocations.refractionTexture, 0);

    this.bindTexture(gl.TEXTURE1, this.framebuffers.reflection.colorTexture);
    gl.uniform1i(program.uniformLocations.reflectionTexture, 1);

    this.bindTexture(gl.TEXTURE2, this.textures.dudvmap);
    gl.uniform1i(program.uniformLocations.dudvTexture, 2);

    this.bindTexture(gl.TEXTURE3, this.textures.normalmap);
    gl.uniform1i(program.uniformLocations.normalMap, 3);

    this.bindTexture(gl.TEXTURE4, this.framebuffers.refraction.depthTexture);
    gl.uniform1i(program.uniformLocations.waterDepthTexture, 4);

    // Draw
    gl.drawElements(gl.TRIANGLES, mesh.indexCount, gl.UNSIGNED_SHORT, 0);
  }

  renderMeshes(clipPlane, mirror) {
    if (!this.state.scenery) return;

    const gl = this.gl;
    const program = this.programs.mesh;
    const mesh = this.meshes.terrain;

    if (!program || !mesh) return;

    gl.useProgram(program);

    // Bind vertex attributes
    this.bindMeshAttributes(program, mesh);

    // Set uniforms
    let viewMatrix = this.getViewMatrix();
    if (mirror) {
      viewMatrix = this.mirrorViewMatrix(viewMatrix);
    }

    const perspectiveMatrix = this.getPerspectiveMatrix();
    const modelMatrix = this.getIdentityMatrix();
    const cameraPos = this.state.camera.position;

    gl.uniformMatrix4fv(
      program.uniformLocations.perspective,
      false,
      perspectiveMatrix,
    );
    gl.uniformMatrix4fv(program.uniformLocations.view, false, viewMatrix);
    gl.uniformMatrix4fv(program.uniformLocations.model, false, modelMatrix);
    gl.uniform3fv(program.uniformLocations.cameraPos, cameraPos);
    gl.uniform4fv(program.uniformLocations.clipPlane, clipPlane);

    // Bind texture
    this.bindTexture(gl.TEXTURE0, this.textures.stone);
    gl.uniform1i(program.uniformLocations.tex, 0);

    // Draw
    gl.drawElements(gl.TRIANGLES, mesh.indexCount, gl.UNSIGNED_SHORT, 0);
  }

  renderDebugViews() {
    const gl = this.gl;
    const program = this.programs.quad;

    if (!program) return;

    gl.useProgram(program);

    // Create a simple quad for rendering textures
    const quadVertices = new Float32Array([
      -1, -1, 0, 0, 1, -1, 1, 0, 1, 1, 1, 1, -1, 1, 0, 1,
    ]);

    const quadIndices = new Uint16Array([0, 1, 2, 0, 2, 3]);

    const quadBuffer = gl.createBuffer();
    gl.bindBuffer(gl.ARRAY_BUFFER, quadBuffer);
    gl.bufferData(gl.ARRAY_BUFFER, quadVertices, gl.STATIC_DRAW);

    const indexBuffer = gl.createBuffer();
    gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer);
    gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, quadIndices, gl.STATIC_DRAW);

    // Enable attributes
    if (program.attribLocations.position !== undefined) {
      gl.enableVertexAttribArray(program.attribLocations.position);
      gl.vertexAttribPointer(
        program.attribLocations.position,
        2,
        gl.FLOAT,
        false,
        16,
        0,
      );
    }

    if (program.attribLocations.texCoords !== undefined) {
      gl.enableVertexAttribArray(program.attribLocations.texCoords);
      gl.vertexAttribPointer(
        program.attribLocations.texCoords,
        2,
        gl.FLOAT,
        false,
        16,
        8,
      );
    }

    // Render reflection preview (top-right corner)
    this.renderDebugQuad(
      program,
      this.framebuffers.reflection.colorTexture,
      this.CANVAS_WIDTH - 100,
      this.CANVAS_HEIGHT - 100,
      100,
      100,
    );

    // Render refraction preview (top-left corner)
    this.renderDebugQuad(
      program,
      this.framebuffers.refraction.colorTexture,
      0,
      this.CANVAS_HEIGHT - 100,
      100,
      100,
    );
  }

  renderDebugQuad(program, texture, x, y, width, height) {
    const gl = this.gl;

    // Set viewport for small preview
    gl.viewport(x, y, width, height);

    // Set orthographic projection for 2D rendering
    const orthoMatrix = this.getOrthoMatrix(-1, 1, -1, 1, -1, 1);
    gl.uniformMatrix4fv(
      program.uniformLocations.perspective,
      false,
      orthoMatrix,
    );
    gl.uniformMatrix4fv(
      program.uniformLocations.view,
      false,
      this.getIdentityMatrix(),
    );
    gl.uniformMatrix4fv(
      program.uniformLocations.model,
      false,
      this.getIdentityMatrix(),
    );

    // Bind texture
    this.bindTexture(gl.TEXTURE0, texture);
    gl.uniform1i(program.uniformLocations.tex, 0);

    // Draw quad
    gl.drawElements(gl.TRIANGLES, 6, gl.UNSIGNED_SHORT, 0);

    // Restore main viewport
    gl.viewport(0, 0, this.CANVAS_WIDTH, this.CANVAS_HEIGHT);
  }

  bindMeshAttributes(program, mesh) {
    const gl = this.gl;

    if (program.attribLocations.position !== undefined) {
      gl.bindBuffer(gl.ARRAY_BUFFER, mesh.vertexBuffer);
      gl.enableVertexAttribArray(program.attribLocations.position);
      gl.vertexAttribPointer(
        program.attribLocations.position,
        3,
        gl.FLOAT,
        false,
        0,
        0,
      );
    }

    if (program.attribLocations.normal !== undefined) {
      gl.bindBuffer(gl.ARRAY_BUFFER, mesh.normalBuffer);
      gl.enableVertexAttribArray(program.attribLocations.normal);
      gl.vertexAttribPointer(
        program.attribLocations.normal,
        3,
        gl.FLOAT,
        false,
        0,
        0,
      );
    }

    if (program.attribLocations.texCoords !== undefined) {
      gl.bindBuffer(gl.ARRAY_BUFFER, mesh.texCoordBuffer);
      gl.enableVertexAttribArray(program.attribLocations.texCoords);
      gl.vertexAttribPointer(
        program.attribLocations.texCoords,
        2,
        gl.FLOAT,
        false,
        0,
        0,
      );
    }

    gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.indexBuffer);
  }

  bindTexture(unit, texture) {
    const gl = this.gl;
    gl.activeTexture(unit);
    gl.bindTexture(gl.TEXTURE_2D, texture);
  }

  // Matrix helper functions
  getPerspectiveMatrix() {
    const fovy = Math.PI / 4; // 45 degrees
    const aspect = this.CANVAS_WIDTH / this.CANVAS_HEIGHT;
    const near = 0.1;
    const far = 1000.0;

    return this.perspective(fovy, aspect, near, far);
  }

  getViewMatrix() {
    if (this.state.camera.viewMatrix) {
      return new Float32Array(this.state.camera.viewMatrix);
    }

    // Fallback view matrix
    const eye = this.state.camera.position;
    const target = [0, 0, 0];
    const up = [0, 1, 0];

    return this.lookAt(eye, target, up);
  }

  getIdentityMatrix() {
    return new Float32Array([1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1]);
  }

  getOrthoMatrix(left, right, bottom, top, near, far) {
    const lr = 1 / (left - right);
    const bt = 1 / (bottom - top);
    const nf = 1 / (near - far);

    return new Float32Array([
      -2 * lr,
      0,
      0,
      0,
      0,
      -2 * bt,
      0,
      0,
      0,
      0,
      2 * nf,
      0,
      (left + right) * lr,
      (top + bottom) * bt,
      (far + near) * nf,
      1,
    ]);
  }

  mirrorViewMatrix(viewMatrix) {
    // Create a mirrored view matrix for reflections
    const mirrored = new Float32Array(viewMatrix);

    // Flip Y-axis components for reflection
    mirrored[1] = -mirrored[1];
    mirrored[5] = -mirrored[5];
    mirrored[9] = -mirrored[9];
    mirrored[13] = -mirrored[13];

    return mirrored;
  }

  perspective(fovy, aspect, near, far) {
    const f = 1.0 / Math.tan(fovy / 2);
    const rangeInv = 1 / (near - far);

    return new Float32Array([
      f / aspect,
      0,
      0,
      0,
      0,
      f,
      0,
      0,
      0,
      0,
      (near + far) * rangeInv,
      -1,
      0,
      0,
      near * far * rangeInv * 2,
      0,
    ]);
  }

  lookAt(eye, center, up) {
    const f = this.normalize(this.subtract(center, eye));
    const s = this.normalize(this.cross(f, up));
    const u = this.cross(s, f);

    return new Float32Array([
      s[0],
      u[0],
      -f[0],
      0,
      s[1],
      u[1],
      -f[1],
      0,
      s[2],
      u[2],
      -f[2],
      0,
      -this.dot(s, eye),
      -this.dot(u, eye),
      this.dot(f, eye),
      1,
    ]);
  }

  // Vector math utilities
  subtract(a, b) {
    return [a[0] - b[0], a[1] - b[1], a[2] - b[2]];
  }

  cross(a, b) {
    return [
      a[1] * b[2] - a[2] * b[1],
      a[2] * b[0] - a[0] * b[2],
      a[0] * b[1] - a[1] * b[0],
    ];
  }

  dot(a, b) {
    return a[0] * b[0] + a[1] * b[1] + a[2] * b[2];
  }

  normalize(v) {
    const len = Math.sqrt(v[0] * v[0] + v[1] * v[1] + v[2] * v[2]);
    if (len > 0) {
      return [v[0] / len, v[1] / len, v[2] / len];
    }
    return [0, 0, 0];
  }
}

// Initialize the application when the page loads
document.addEventListener("DOMContentLoaded", () => {
  new WebGLWaterApp();
});
